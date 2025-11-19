package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader переводит HTTP в WebSocket.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// в проде здесь надо проверять Origin.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// echoHandler апгрейдит соединение и эхо‑ит сообщения.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	// апгрейд HTTP -> WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// читаем сообщение от клиента
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("recv: %s\n", msg)

		// отправляем назад то же сообщение
		if err := conn.WriteMessage(mt, msg); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", echoHandler)

	log.Println("WebSocket server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
