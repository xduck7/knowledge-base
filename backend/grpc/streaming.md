# gRPC streaming 

> В gRPC есть четыре типа RPC: unary, server‑streaming, client‑streaming и bidirectional‑streaming.  
> Стриминг даёт возможность передавать последовательности сообщений поверх одного HTTP/2‑стрима с сохранением порядка и управлением потоком.

---

## Типы стриминга в gRPC

Кратко:

- **Unary**  
  Классический запрос‑ответ
- **Server‑streaming**  
  Клиент отправляет один запрос, сервер возвращает поток сообщений
- **Client‑streaming**  
  Клиент отправляет поток сообщений, получает один ответ
- **Bidirectional‑streaming**  
  Обе стороны могут одновременно читать и писать в поток

В `.proto` тип стриминга задаётся ключевым словом `stream` перед типом запроса/ответа:

```protobuf
service ChatService {
// unary
rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);

// server streaming
rpc ListMessages(ListMessagesRequest) returns (stream Message);

// client streaming
rpc UploadEvents(stream Event) returns (UploadSummary);

// bidi streaming
rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
```

---

## Server‑streaming

### Схема

- Клиент → один запрос.
- Сервер → последовательность сообщений в ответ.
- Поток закрывается, когда сервер закончил отправку или произошла ошибка.

Примеры кейсов:

- выдача большого списка (логи, записи, результаты поиска);
- подписка на события за ограниченный промежуток времени;
- передача файла чанками от сервера к клиенту.

### Сервер (Go)

```go
func (s *chatServer) ListMessages(
req *chatv1.ListMessagesRequest,
stream chatv1.ChatService_ListMessagesServer,
) error {
messages, err := s.repo.List(stream.Context(), req.GetChatId())
if err != nil {
return status.Errorf(codes.Internal, "list messages: %v", err)
}

    for _, m := range messages {
        if err := stream.Send(&chatv1.Message{
            Id:      m.ID,
            ChatId:  m.ChatID,
            Text:    m.Text,
            Created: timestamppb.New(m.CreatedAt),
        }); err != nil {
            return err // при ошибке отправки выходим
        }
    }

    return nil
}
```

### Клиент (Go)

```go
stream, err := client.ListMessages(ctx, &chatv1.ListMessagesRequest{ChatId: chatID})
if err != nil {
log.Fatalf("ListMessages error: %v", err)
}

for {
msg, err := stream.Recv()
if errors.Is(err, io.EOF) {
break // сервер закончил стрим
}
if err != nil {
log.Fatalf("recv error: %v", err)
}

    log.Printf("msg: %s", msg.GetText())
}
```

---

## Client‑streaming

### Схема

- Клиент → поток сообщений.
- Сервер → один итоговый ответ (после закрытия потока клиентом).

Примеры кейсов:

- загрузка данных чанками (файл, лог, батч событий);
- агрегирование событий на сервере (метрики, логи, трекинг);
- отправка большого набора данных, для которого нужен только summary.

### Клиент (Go)

```go
stream, err := client.UploadEvents(ctx)
if err != nil {
log.Fatalf("UploadEvents error: %v", err)
}

for _, e := range events {
if err := stream.Send(&eventsv1.Event{
Id:   e.ID,
Type: e.Type,
Data: e.Data,
}); err != nil {
log.Fatalf("send error: %v", err)
}
}

// важно: закрыть поток записи
summary, err := stream.CloseAndRecv()
if err != nil {
log.Fatalf("CloseAndRecv error: %v", err)
}

log.Printf("uploaded: %d events", summary.GetCount())
```

### Сервер (Go)

```go
func (s *eventServer) UploadEvents(
stream eventsv1.EventService_UploadEventsServer,
) error {
var count int64

    for {
        ev, err := stream.Recv()
        if errors.Is(err, io.EOF) {
            // клиент закончил отправку — возвращаем summary
            return stream.SendAndClose(&eventsv1.UploadSummary{
                Count: count,
            })
        }
        if err != nil {
            return status.Errorf(codes.Internal, "recv: %v", err)
        }

        if err := s.repo.Save(stream.Context(), ev); err != nil {
            return status.Errorf(codes.Internal, "save: %v", err)
        }

        count++
    }
}
```

---

## Bidirectional‑streaming

### Схема

- Клиент и сервер могут **одновременно**:
  - отправлять сообщения в поток;
  - читать сообщения из потока.
- Порядок сообщений **в каждом направлении** сохраняется, но стороны не обязаны чередовать «отправил — прочитал».

Примеры кейсов:

- чаты и уведомления в реальном времени;
- online‑игры и синхронизация состояний;
- стриминг телеметрии с ответными командами от сервера.

### Клиент (Go, упрощённо)

```go
stream, err := client.Chat(ctx)
if err != nil {
log.Fatalf("Chat error: %v", err)
}

// горутина для приёма
go func() {
    for {
    msg, err := stream.Recv()
        if errors.Is(err, io.EOF) {
            return
        }
        if err != nil {
            log.Printf("recv error: %v", err)
            return
        }
        log.Printf("incoming: %s", msg.GetText())
    }
}()

// отправка
for _, text := range []string{"hi", "how are you?"} {
    if err := stream.Send(&chatv1.ChatMessage{Text: text}); err != nil {
        log.Fatalf("send error: %v", err)
    }
}

if err := stream.CloseSend(); err != nil {
    log.Fatalf("CloseSend error: %v", err)
}
```

---

## Поток, контекст и ошибки

Общие моменты для всех видов стриминга:

- **Контекст**:
  - `stream.Context()` на сервере даёт доступ к `deadline`, `cancellation`, метаданным;
  - при отмене/таймауте клиента `Context` на сервере тоже отменяется.
- **Ошибки**:
  - `io.EOF` на клиенте/сервере означает «другая сторона закрыла поток» (нормальное завершение);
  - любые другие ошибки → нужно корректно завершить RPC с gRPC‑статусом (через `status.Errorf`).
- **Закрытие**:
  - client‑streaming: клиент обязан вызвать `CloseAndRecv` (или `CloseSend` для bidi);
  - server‑streaming: сервер завершает метод, клиент видит `io.EOF`.

---

## Когда стоит (и не стоит) использовать стриминг

**Подходит, когда:**

- много однотипных сообщений (лог/метрика/события);
- большой объём данных, который не хочется держать целиком в памяти;
- нужен realtime‑канал (подписки, уведомления, игры, чаты).

**Не подходит, когда:**

- запрос‑ответ строго один‑к‑одному и payload маленький → лучше unary;
- нужна простая HTTP‑интеграция с внешними клиентами (там проще REST/WebSocket).

Частая ошибка — делать streaming там, где обычный unary идеально подходит. Стриминг усложняет код (жизненный цикл, ошибки, конкуренция), поэтому его стоит включать только тогда, когда реально нужен.

---

## Практические советы

- Всегда задавай таймауты на стриминговые вызовы (особенно long‑lived).
- Логируй открытие и закрытие стримов, duration, количество сообщений.
- Бережно относись к backpressure:
  - не читай/не пиши бесконечно без пауз;
  - учитывай, что в Go `Recv` блокируется, пока нет данных.
- Для bidi‑streaming аккуратно работай с горутинами:
  - одна горутина читает, другая пишет;
  - грамотно обрабатывай закрытие обеих сторон.
