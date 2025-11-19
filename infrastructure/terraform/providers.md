# Провайдеры Terraform

## Что такое провайдер

Провайдер в Terraform — это плагин, который отвечает за взаимодействие с конкретной платформой или сервисом.  
Он знает, как создавать, читать, обновлять и удалять ресурсы на стороне облака, SaaS или локальной инфраструктуры.

Примеры провайдеров:

- `aws` — Amazon Web Services  
- `google` — Google Cloud Platform  
- `azurerm` — Microsoft Azure  
- `kubernetes` — Kubernetes cluster  
- `vault` — HashiCorp Vault  
- `random` — генерация случайных значений  

Провайдеры определяют набор ресурсов и data source'ов, с которыми Terraform работает.

---

## Подключение провайдера

Чтобы использовать провайдер, нужно:

1. Задать его в блоке `terraform { required_providers { ... } }`  
2. Настроить подключение через блок `provider`

Пример для AWS:

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region  = "us-east-1"
  profile = "default"
}
````

* `source` — откуда брать провайдер (Terraform Registry).
* `version` — версия плагина, рекомендуется фиксировать версию.
* В блоке `provider` указываются параметры подключения (регион, учётные данные, endpoint).

---

## Несколько провайдеров одного типа

Можно использовать несколько провайдеров одного типа для разных регионов или аккаунтов, используя алиасы:

```hcl
provider "aws" {
  alias  = "us_east"
  region = "us-east-1"
}

provider "aws" {
  alias  = "eu_west"
  region = "eu-west-1"
}

resource "aws_s3_bucket" "bucket_us" {
  provider = aws.us_east
  bucket   = "my-us-bucket"
}

resource "aws_s3_bucket" "bucket_eu" {
  provider = aws.eu_west
  bucket   = "my-eu-bucket"
}
```

* `alias` — уникальное имя для провайдера.
* `provider = aws.us_east` в ресурсе указывает, какой провайдер использовать.

---

## Провайдеры и переменные окружения

Многие провайдеры могут брать параметры подключения из переменных окружения:

* AWS: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`
* Google Cloud: `GOOGLE_CREDENTIALS`, `GOOGLE_PROJECT`, `GOOGLE_REGION`
* Azure: `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`, `ARM_SUBSCRIPTION_ID`, `ARM_TENANT_ID`

Пример для AWS без явного `provider`:

```sh
export AWS_ACCESS_KEY_ID="AKIA..."
export AWS_SECRET_ACCESS_KEY="..."
export AWS_REGION="us-east-1"
terraform plan
```

Terraform автоматически подхватит эти значения.

---

## Data sources провайдера

Провайдеры позволяют не только создавать ресурсы, но и получать существующие через `data`:

```hcl
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Ubuntu official
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
}
```

* Data source не создаёт объект, а только получает информацию.
* Полезно для динамического использования ID ресурсов, образов, сетей и т.п.

---

## Рекомендации по работе с провайдерами

* Всегда фиксировать версии провайдеров, чтобы конфигурация была предсказуемой.
* Использовать алиасы для нескольких аккаунтов/регионов.
* Стараться использовать переменные для параметров подключения, чтобы не хранить секреты в коде.
* Проверять документацию провайдера для допустимых параметров ресурсов и data source.

---

## Итог

* Провайдер — мост между Terraform и внешними системами.
* Определяется в `terraform.required_providers` и настраивается через `provider`.
* Поддерживает ресурсы и data source.
* Для нескольких аккаунтов/регионов использовать алиасы.
* Переменные окружения упрощают безопасное хранение учетных данных.
