# Структура Helm‑чарта

Helm‑чарт — это директория с предсказуемой структурой, в которой лежат метаданные, значения по умолчанию, шаблоны манифестов и, при необходимости, зависимости и тесты.  
Понимание этой структуры помогает быстро ориентироваться в чужих чартах и писать свои так, чтобы их было удобно сопровождать.

Типовая структура:

```sh
mychart/
├── Chart.yaml
├── values.yaml
├── charts/
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── _helpers.tpl
│   ├── NOTES.txt
│   └── tests/
│       └── test-connection.yaml
└── .helmignore
```

## Chart.yaml

`Chart.yaml` содержит метаданные чарта:

- `name` — имя чарта;  
- `version` — версия самого чарта (не приложения);  
- `appVersion` — версия приложения, которое деплоится;  
- `description` — краткое описание;  
- `type` — `application` или `library`;  
- `dependencies` — список зависимостей (subchart’ы).

Пример:

```c
apiVersion: v2
name: myapp
description: My awesome backend service
type: application
version: 0.1.0
appVersion: "1.0.3"

dependencies:
  - name: redis
    version: 17.3.3
    repository: https://charts.bitnami.com/bitnami
```

## values.yaml

`values.yaml` — файл со значениями по умолчанию, которые подставляются в шаблоны.  
Идея: в шаблонах как можно меньше «жёстко пришитых» значений, всё выносится в values.

Пример:

```c
replicaCount: 2

image:
  repository: my-registry/myapp
  tag: "1.0.3"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "500m"
    memory: "512Mi"
```

## templates/

Директория `templates/` содержит шаблоны Kubernetes‑манифестов (с Go‑template синтаксисом Helm):

- `deployment.yaml` — шаблон Deployment;  
- `service.yaml` — Service;  
- `ingress.yaml` — Ingress;  
- другие ресурсы (ConfigMap, Secret, HPA, Job и т.д.);  
- `_helpers.tpl` — вспомогательные функции/шаблоны;  
- `NOTES.txt` — текст, который Helm выводит после `install/upgrade` (подсказки по доступу, командам и т.п.);  
- `tests/` — тестовые ресурсы для `helm test`.

Пример заголовка `deployment.yaml`:

```c
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "myapp.fullname" . }}
  labels:
    {{- include "myapp.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "myapp.selectorLabels" . | nindent 6 }}
```

## charts/

Папка `charts/` содержит зависимости (subchart’ы):

- если чарт объявляет зависимости в `Chart.yaml`, при `helm dependency update` они скачиваются сюда;  
- такой чарт можно потом использовать как «зонтик» над несколькими компонентами.

## .helmignore

`.helmignore` задаёт, какие файлы игнорировать при упаковке чарта (`helm package`):

```sh
.git/
*.swp
*.tmp
README.md
tests/
```

Это уменьшает размер пакета и исключает из него мусор.

## tests/

Тесты обычно лежат в `templates/tests/` и представляют собой Kubernetes‑ресурсы с аннотацией `helm.sh/hook: test`.  
Например, `test-connection.yaml`, который запускает pod и пытается достучаться до сервиса.

```c
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "myapp.fullname" . }}-test-connection"
  annotations:
    "helm.sh/hook": test
spec:
  containers:
    - name: curl
      image: curlimages/curl
      args: ["curl", "{{ include "myapp.fullname" . }}:80/healthz"]
  restartPolicy: Never
```

# Шаблоны и функции в Helm

## Базовый синтаксис шаблонов

Helm использует Go templates: выражения обрамляются в `{{ ... }}` и могут появляться в любом месте YAML.  
Шаблон получает на вход объект `.`, в котором лежит контекст чарта: `.Values`, `.Chart`, `.Release`, `.Capabilities` и т.д.

Простейший пример:

```c
metadata:
  name: {{ .Release.Name }}-backend
```

Здесь `{{ .Release.Name }}` будет заменён на имя релиза (`helm install myapp ...` → `myapp-backend`).

## Обращение к values

Чаще всего шаблоны обращаются к значениям из `values.yaml`:

```c
spec:
  replicas: {{ .Values.replicaCount }}
  template:
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
```

Если какого‑то ключа нет, по умолчанию подставится пустое значение (но лучше явно задавать всё необходимое в `values.yaml`).

## Функции и pipe

Helm поверх Go templates подтягивает набор функций (из Sprig и своих):

- `default` — значение по умолчанию;  
- `upper`, `lower`, `replace`, `trunc` и т.д.;  
- `toYaml`, `nindent` — удобство для форматирования;  
- и многие другие.

Функции комбинируются через pipe, как в shell:

```c
metadata:
  labels:
    app.kubernetes.io/name: {{ .Chart.Name | quote }}
    app.kubernetes.io/instance: {{ .Release.Name | quote }}
```

Пример с `default` и `toYaml`:

```js
resources:
{{- toYaml (default .Values.resources .Values.backend.resources) | nindent 2 }}
```

Здесь:

- сначала выбирается либо `.Values.backend.resources`, либо `.Values.resources`, если первый пустой;  
- `toYaml` превращает map в YAML;  
- `nindent 2` смещает текст на два пробела.

## Условные конструкции

Условный рендеринг делается через `if/else`:

```sh
{{- if .Values.ingress.enabled }}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ include "myapp.fullname" . }}
spec:
  ...
{{- end }}
```

Такой ресурс попадёт в рендер только если `ingress.enabled: true` в values.  
Можно вкладывать условия, использовать `and`, `or`, `not`:

```c
{{- if and .Values.metrics.enabled .Values.serviceMonitor.enabled }}
# ServiceMonitor ресурс...
{{- end }}
```

## Циклы range

Для повторяющихся сущностей используют `range`:

```sh
env:
  {{- range $key, $value := .Values.env }}
  - name: {{ $key }}
    value: {{ $value | quote }}
  {{- end }}
```

При `env`:

```sh
env:
  LOG_LEVEL: info
  FEATURE_FLAG_X: "true"
```

получится список переменных окружения в контейнере.

Аналогично можно рендерить несколько портов, volume’ы, sidecar‑контейнеры и т.д.

## include / define и _helpers.tpl

Файл `_helpers.tpl` содержит переиспользуемые куски шаблонов (partial’ы), объявленные через `define`:

```sh
{{- define "myapp.name" -}}
{{ .Chart.Name }}
{{- end }}

{{- define "myapp.fullname" -}}
{{ printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "myapp.labels" -}}
app.kubernetes.io/name: {{ include "myapp.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
```

Использование:

```sh
metadata:
  name: {{ include "myapp.fullname" . }}
  labels:
    {{- include "myapp.labels" . | nindent 4 }}
```

Преимущества:

- единые правила именования ресурсов;  
- единый набор стандартных label’ов/annotation’ов;  
- меньше дублирования в `deployment.yaml`, `service.yaml`, `ingress.yaml` и т.д.

## NOTES.txt

`templates/NOTES.txt` — не шаблон манифеста, а текст, который выводит Helm после `install`/`upgrade`.  
В нём можно подсказать пользователю:

- как проверить статус;  
- как получить URL сервиса;  
- как залогиниться и т.п.

Пример:

```
1. Get the application URL by running:

  export POD_NAME=$(kubectl get pods --namespace {{ .Release.Namespace }} -l "app={{ include "myapp.name" . }}" -o jsonpath="{.items.metadata.name}")
  kubectl port-forward --namespace {{ .Release.Namespace }} $POD_NAME 8080:80
```

NOTES тоже можно шаблонизировать через `{{ }}`.

# Работа с values.yaml

`values.yaml` — ключевое место, где настраивается поведение чарта без изменения шаблонов.  
Через него описываются ресурсы, ингрессы, переменные окружения, подключение к БД, включение/выключение фич и многое другое.

## Значения по умолчанию

В корневом `values.yaml` задаются sane defaults:

```sh
replicaCount: 2

image:
  repository: my-registry/myapp
  tag: "1.0.0"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  className: ""
  hosts: []
  tls: []

resources: {}
env: {}
```

Шаблоны используют эти значения, но в любой установке их можно переопределить.

## Переопределение через -f и --set

При установке/обновлении чарта значения можно изменять:

- через дополнительные файлы: `-f custom-values.yaml`;  
- через CLI: `--set key=value` (и `--set-string`, `--set-json`).

Примеры:

```sh
# Базовая установка с доп. values
helm install myapp ./chart -f values.prod.yaml

# Переопределение отдельных параметров
helm upgrade myapp ./chart \
  --set replicaCount=3 \
  --set image.tag=1.0.5 \
  --set ingress.enabled=true
```

Рекомендации:

- для прод/стейджа/дева лучше делать отдельные файлы `values.prod.yaml`, `values.stage.yaml`, `values.dev.yaml`;  
- `--set` использовать для быстрых экспериментов или в CI, когда меняется пара значений.

## Организация per‑env values

Один чарт может использоваться во всех окружениях, различаются только values:

```sh
myapp/
  Chart.yaml
  values.yaml            # значения по умолчанию
  values.dev.yaml        # dev-окружение
  values.stage.yaml      # stage
  values.prod.yaml       # prod
```

Типичный запуск:

```sh
helm install myapp-dev ./myapp -f values.dev.yaml
helm install myapp-stage ./myapp -f values.stage.yaml
helm install myapp-prod ./myapp -f values.prod.yaml
```

Часто:

- `values.yaml` содержит разумные дефолты (минимальные ресурсы, выключенные ингрессы/мониторинг);  
- в env‑специфичных файлах только отличия: репликасы, ингрессы, ресурсы, URL внешних сервисов.

## Типичные паттерны: resources, ingress, env, secrets

### resources

Блок ресурсов обычно делается один и тот же:

```sh
resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "500m"
    memory: "512Mi"
```

В шаблоне:

```sh
resources:
  {{- toYaml .Values.resources | nindent 2 }}
```

На проде в `values.prod.yaml` можно задать более высокие лимиты, на dev — пониже.

### ingress

Гибкий ингресс‑блок:

```sh
ingress:
  enabled: true
  className: nginx
  hosts:
    - host: myapp.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: myapp-tls
      hosts:
        - myapp.example.com
```

В шаблоне:

```sh
{{- if .Values.ingress.enabled }}
apiVersion: networking.k8s.io/v1
kind: Ingress
...
{{- end }}
```

Так можно включать/выключать ingress по окружениям и менять хосты/классы.

### env (переменные окружения)

Распространённый паттерн — map в values, рендер через `range`:

```sh
env:
  LOG_LEVEL: info
  FEATURE_FLAG_X: "true"
```

```sh
env:
  {{- range $key, $value := .Values.env }}
  - name: {{ $key }}
    value: {{ $value | quote }}
  {{- end }}
```

Для чувствительных данных вместо value используют `valueFrom` и Kubernetes Secret, но сами секреты лучше не держать в открытом `values.yaml`.

### secrets и конфиги

Обычно:

- сами секреты живут отдельно (SealedSecrets, SOPS, ExternalSecrets);  
- в values задаются только ссылки/имена:

```sh
secret:
  name: myapp-secret
  key: DATABASE_URL
```

В шаблоне Deployment:

```sh
env:
  - name: DATABASE_URL
    valueFrom:
      secretKeyRef:
        name: {{ .Values.secret.name }}
        key: {{ .Values.secret.key }}
```

Так не приходится таскать чувствительные данные в helm‑values, но чарт остаётся конфигурируемым.