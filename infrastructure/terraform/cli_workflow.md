# Terraform CLI workflow

## Базовый цикл работы

Terraform CLI крутится вокруг трёх основных шагов: инициализация проекта, просмотр плана изменений и применение этих изменений.  
Обычно это выглядит как последовательность `terraform init → terraform plan → terraform apply`, плюс иногда `destroy` для удаления инфраструктуры.  

---

## terraform init

`terraform init` подготавливает директорию с конфигами:

- скачивает и инициализирует провайдеры;  
- настраивает backend (где хранится state);  
- генерирует служебные файлы в `.terraform/`.  

Вызывается:

```bash
terraform init
```

Его запускают:

- один раз при создании проекта/модуля;  
- при изменении конфигурации backend’а;  
- при добавлении/обновлении провайдеров.  

---

## terraform fmt и validate

Перед планом удобно прогнать форматирование и базовую валидацию:

- `terraform fmt` — приводит `.tf` файлы к каноничному стилю;  
- `terraform validate` — проверяет синтаксис и базовую корректность конфигурации.  

Пример:

```bash
terraform fmt
terraform validate
```

Это дешёвые проверки, которые хорошо заходят в pre-commit/CI.

---

## terraform plan

`terraform plan` показывает, **что** Terraform собирается сделать, не внося изменений.  
Он сравнивает желаемое состояние (конфиги) с текущим (state + реальность) и строит план действий.

```bash
terraform plan
```

Можно:

- вывести план в файл:  

  ```bash
  terraform plan -out=tfplan.bin
  ```

- ограничить подмножество ресурсов через `-target`, но это стоит использовать осторожно.  

План — главный инструмент ревью изменений: смотреть, какие ресурсы будут созданы/изменены/удалены.

---

## terraform apply

`terraform apply` применяет изменения:

```sh
terraform apply
```

Либо по подготовленному плану:

```sh
terraform plan  -out=tfplan.bin
terraform apply tfplan.bin
```

На проде часто используют именно двухшаговый вариант:  
- в CI генерируют и сохраняют план;  
- на ручном шаге `apply` применяют ровно его.  

---

## terraform destroy

`terraform destroy` делает инверсию плана — удаляет все ресурсы, описанные в конфигурации/state:

```bash
terraform destroy
```

Иногда:

- используют для временных стендов и окружений для тестов;  
- комбинируют с `-target` для выборочного удаления.  

Нужно помнить, что destroy работает относительно state: если часть ресурсов уже удалена руками, Terraform всё равно попытается привести состояние к «ничего нет».

---

## Полезные доп. команды

- `terraform show` — показывает текущее состояние (state) в человекочитаемом виде.  
- `terraform state list` / `show <addr>` — работа с конкретными ресурсами в state.  
- `terraform import` — подтягивает уже существующий ресурс в управление Terraform.  
- `terraform output` — печатает outputs root‑модуля (часто используют в CI/скриптах).  

---

## Пример удобного sh‑скрипта

Ниже пример простого скрипта, который даёт единый интерфейс `./tf.sh up` и `./tf.sh down`:

```bash
#!/usr/bin/env bash
set -euo pipefail

CMD="${1:-}"

usage() {
  echo "Usage: $0 {up|down|plan|apply|destroy}"
  exit 1
}

info()  { echo "[INFO]  $*"; }
error() { echo "[ERROR] $*" >&2; }

tf_init() {
  info "Initializing Terraform..."
  terraform init -input=false
}

tf_fmt_validate() {
  info "Formatting and validating..."
  terraform fmt -recursive
  terraform validate
}

tf_plan() {
  tf_init
  tf_fmt_validate
  info "Planning..."
  terraform plan -out=tfplan.bin
}

tf_apply() {
  if [[ ! -f tfplan.bin ]]; then
    tf_plan
  fi
  info "Applying plan..."
  terraform apply -input=false tfplan.bin
}

tf_destroy() {
  tf_init
  info "Destroying..."
  terraform destroy
}

case "${CMD}" in
  up)
    # Полный цикл: init + fmt/validate + plan + apply
    tf_apply
    ;;
  down)
    # Удаление инфраструктуры
    tf_destroy
    ;;
  plan)
    tf_plan
    ;;
  apply)
    tf_apply
    ;;
  destroy)
    tf_destroy
    ;;
  *)
    usage
    ;;
esac
```

Как использовать:

```
chmod +x tf.sh

# Развернуть/обновить инфраструктуру одной командой
./tf.sh up

# Снести всё
./tf.sh down
```
