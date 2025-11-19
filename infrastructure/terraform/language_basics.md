# Основы языка Terraform (HCL)

## Что такое HCL

HCL (HashiCorp Configuration Language) — декларативный язык конфигураций, используемый в Terraform для описания инфраструктуры.  
Цель HCL — быть читаемым для человека и удобным для автоматического парсинга.

---

## Структура `.tf` файлов

Terraform конфигурации обычно хранятся в файлах с расширением `.tf`.  
Типовая структура проекта:

```

project/
├── main.tf       # основные ресурсы
├── variables.tf  # описание переменных
├── outputs.tf    # выходные значения
├── providers.tf  # подключение провайдеров
└── terraform.tfvars  # значения переменных

````

Terraform автоматически объединяет все `.tf` файлы в одной директории в единую конфигурацию.

---

## Основные блоки конфигурации

### terraform

Блок `terraform` задаёт глобальные настройки Terraform, backend, required providers:

```hcl
terraform {
  required_version = ">= 1.5.0"

  backend "s3" {
    bucket = "my-terraform-state"
    key    = "prod/terraform.tfstate"
    region = "us-east-1"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
````

---

### provider

Блок `provider` указывает Terraform, как подключаться к облачному сервису или API:

```hcl
provider "aws" {
  region = "us-east-1"
}
```

Можно задавать несколько провайдеров для разных облаков или регионов.

---

### resource

Блок `resource` создаёт конкретный объект инфраструктуры:

```hcl
resource "aws_s3_bucket" "my_bucket" {
  bucket = "my-example-bucket"
  acl    = "private"
  tags = {
    Environment = "prod"
    Team        = "backend"
  }
}
```

* Первый параметр (`aws_s3_bucket`) — тип ресурса.
* Второй (`my_bucket`) — локальное имя в конфигурации.
* Внутри блока задаются параметры ресурса.

---

### data

Блок `data` позволяет получить информацию о существующих объектах, чтобы использовать её в конфигурации:

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

---

## Базовые типы данных

Terraform поддерживает следующие типы:

* `string` — текстовая строка:

  ```hcl
  variable "region" {
    type    = string
    default = "us-east-1"
  }
  ```
* `number` — число:

  ```hcl
  variable "instance_count" {
    type    = number
    default = 3
  }
  ```
* `bool` — логическое значение:

  ```hcl
  variable "enable_feature" {
    type    = bool
    default = true
  }
  ```
* `list` — упорядоченный массив:

  ```hcl
  variable "zones" {
    type    = list(string)
    default = ["us-east-1a", "us-east-1b"]
  }
  ```
* `map` — ключ-значение:

  ```hcl
  variable "tags" {
    type = map(string)
    default = {
      Environment = "prod"
      Team        = "backend"
    }
  }
  ```

---

## Итог

* HCL — декларативный язык для описания инфраструктуры.
* Основные блоки: `terraform`, `provider`, `resource`, `data`.
* Конфигурации хранятся в `.tf` файлах и объединяются автоматически.
* Базовые типы данных: string, number, bool, list, map.
