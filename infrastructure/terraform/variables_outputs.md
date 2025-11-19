## Зачем нужны variables и outputs

Переменные и outputs — это «входы» и «выходы» Terraform‑конфигураций.  
Переменные позволяют параметризовать инфраструктуру (регион, размеры, фичи), а outputs — удобно забирать важные значения (IP, ID, URL) из уже созданных ресурсов и передавать их дальше.  

---

## Входные переменные (input variables)

### Что это

Input‑переменные описывают параметры, которые конфигурация ожидает «снаружи».  
Они делают код переиспользуемым: один и тот же модуль можно запустить в разных окружениях, меняя только значения переменных.  

Переменные обычно объявляют в `variables.tf` (не обязательно, но так принято).  

### Объявление переменных

Базовый синтаксис:

```tf
variable "region" {
  description = "Cloud region"
  type        = string
  default     = "eu-central-1"
}

variable "instance_count" {
  description = "Number of app instances"
  type        = number
  default     = 2
}

variable "tags" {
  description = "Common resource tags"
  type        = map(string)
  default = {
    project = "my-service"
    env     = "dev"
  }
}
```

Ключевые поля:

- `description` — документация (очень помогает в командах);  
- `type` — тип (`string`, `number`, `bool`, `list`, `map`, `set`, `object`, `tuple` и т.д.);  
- `default` — значение по умолчанию (если нет, переменная становится обязательной).  

### Использование переменных

Переменная доступна как `var.<name>`:

```tf
provider "aws" {
  region = var.region
}

resource "aws_instance" "app" {
  count = var.instance_count

  tags = merge(
    var.tags,
    {
      "component" = "app"
    }
  )
}
```

---

## Как задавать значения переменных

Terraform поддерживает несколько источников значений, с приоритетами.

### Файлы .tfvars

Самый удобный способ — `*.tfvars` / `*.tfvars.json`:

```tf
# dev.tfvars
region         = "eu-central-1"
instance_count = 1
```

Запуск:

```bash
terraform apply -var-file="dev.tfvars"
```

Часто используют:

- `dev.tfvars`, `stage.tfvars`, `prod.tfvars`;  
- `terraform.tfvars` — подхватывается автоматически, если есть.  

### CLI‑флаг -var

Для точечных переопределений:

```bash
terraform apply -var="region=eu-west-1" -var="instance_count=3"
```

Удобно для CI/скриптов, но плохо читаемо, если параметров много.  

### Переменные окружения

Любая переменная вида `TF_VAR_<name>` воспринимается как input‑переменная:

```bash
export TF_VAR_region=eu-north-1
terraform apply
```

Это удобно для секретов (вместо tfvars в Git), но лучше в проде комбинировать с секрет‑менеджерами.  

### Интерактивный ввод

Если у переменной нет `default` и значение нигде не задано, Terraform при `plan/apply` спросит его интерактивно.  
В продовых пайплайнах это обычно отключают, чтобы всё было явно в tfvars/окружении.

---

## Типы и валидация

### Типы

Тип задают через `type`:

```tf
variable "allowed_cidrs" {
  type        = list(string)
  description = "CIDR blocks allowed to access the app"
}

variable "config" {
  type = object({
    cpu    = number
    memory = number
  })
}
```

Статическая типизация помогает ловить ошибки до того, как Terraform полезет в провайдеры.  

### Валидация (validation)

Можно навесить дополнительные проверки:

```tf
variable "env" {
  type = string

  validation {
    condition     = contains(["dev", "stage", "prod"], var.env)
    error_message = "env must be one of dev, stage, prod."
  }
}
```

---

## Выходные значения (outputs)

### Что это

Outputs — это именованные значения, которые Terraform «выдаёт наружу» после `apply`.  
Они нужны, чтобы:

- быстро посмотреть важные данные (IP, URL, ID);  
- передать значения из одного модуля другому;  
- отдать результаты во внешние системы (CI, скрипты).  

Обычно их кладут в `outputs.tf` (для читаемости).  

### Объявление output’ов

Простой пример:

```tf
output "app_public_ip" {
  description = "Public IP of the app instance"
  value       = aws_instance.app.public_ip
}

output "app_url" {
  description = "Application URL"
  value       = "https://${aws_lb.app.dns_name}"
}
```

`terraform apply` в конце покажет:

```bash
Outputs:

app_public_ip = "203.0.113.10"
app_url       = "https://my-app-lb-1234.elb.amazonaws.com"
```

### Чтение output’ов

После успешного применения:

```bash
terraform output           # все outputs
terraform output app_url   # один конкретный
```

В скриптах можно запросить JSON:

```bash
terraform output -json
terraform output -json app_url
```

---

## Outputs и модули

### Передача значений между модулями

Outputs — основной механизм, через который модуль даёт наружу свои «результаты».  
Внешний код может использовать их как обычные значения.

Модуль `vpc`:

```tf
# modules/vpc/outputs.tf
output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_subnet_ids" {
  value = aws_subnet.public[*].id
}
```

Root‑модуль:

```tf
module "vpc" {
  source = "./modules/vpc"
  # ...
}

resource "aws_instance" "app" {
  subnet_id = module.vpc.public_subnet_ids
  # ...
}
```

Здесь:

- `module.vpc.vpc_id` и `module.vpc.public_subnet_ids` — ссылки на outputs модуля `vpc`;  
- граф зависимостей автоматически понимает, что `aws_instance.app` зависит от модуля `vpc`.  

---

## Структура файлов: variables.tf и outputs.tf

На практике конфигурацию обычно делят примерно так:

- `main.tf` — сами ресурсы;  
- `providers.tf` — провайдеры;  
- `variables.tf` — объявления переменных;  
- `terraform.tfvars` / `*.tfvars` — значения для конкретного окружения;  
- `outputs.tf` — outputs.  

Это не требование Terraform, а устоявшийся стиль: так проще читать и ревьюить конфигурации.
