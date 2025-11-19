## Рекомендации по структуре чарта

Структура чарта должна быть предсказуемой. 

Шаблоны находятся в `templates/`, 

зависимости — в `charts/`, 

значения — в `values.yaml`. 

Именование ресурсов идёт через функции Helm, чтобы имена были стабильными. Labels и annotations следуют Kubernetes-стандартам и позволяют трекать релиз и версию приложения. Один чарт описывает одну сущность. Линт и тесты не допускают разъезда шаблонов.

Пример минимального шаблона Deployment:

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "service.fullname" . }}
  labels:
    app.kubernetes.io/name: {{ include "service.name" . }}
    app.kubernetes.io/instance: {{ .Release.Name }}
    app.kubernetes.io/version: {{ .Chart.Version }}
    app.kubernetes.io/managed-by: helm
spec:
  replicas: {{ .Values.replicas }}
  selector:
    matchLabels:
      app.kubernetes.io/name: {{ include "service.name" . }}
  template:
    metadata:
      labels:
        app.kubernetes.io/name: {{ include "service.name" . }}
    spec:
      containers:
        - name: app
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          ports:
            - name: http
              containerPort: {{ .Values.service.port }}
```

Пример `_helpers.tpl`:

```tpl
{{- define "service.name" -}}
{{ .Chart.Name }}
{{- end -}}

{{- define "service.fullname" -}}
{{ .Release.Name }}-{{ .Chart.Name }}
{{- end -}}
```

---

### Работа с values и гигиена конфигурации

`values.yaml` содержит только sane defaults. Параметры для окружений вынесены в отдельные файлы. Иерархия значений формализована: image, service, resources, probes.

Пример базового `values.yaml`:

```yaml
replicas: 1

image:
  repository: my-registry/app
  tag: "latest"
  pullPolicy: IfNotPresent

service:
  port: 8080

resources:
  limits:
    cpu: "500m"
    memory: "256Mi"
  requests:
    cpu: "100m"
    memory: "128Mi"

probes:
  liveness:
    path: /live
    port: 8080
  readiness:
    path: /ready
    port: 8080
```

Пример файла для прод-окружения:

```yaml
# values-prod.yaml
replicas: 3

image:
  tag: "1.4.7"

resources:
  limits:
    cpu: "1"
    memory: "512Mi"
  requests:
    cpu: "250m"
    memory: "256Mi"
```

---

### Безопасность и управление секретами

Секреты не хранятся в открытом виде. Используются Sealed Secrets, SOPS или External Secrets. RBAC — минимальный. Defaults — безопасные.

Пример SealedSecret:

```yaml
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  name: db-password
spec:
  encryptedData:
    password: AgBvJf1h1...
```

Пример интеграции External Secrets:

```yaml
# templates/external-secret.yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: db-creds
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: vault-backend
    kind: ClusterSecretStore
  target:
    name: db-creds
  data:
    - secretKey: password
      remoteRef:
        key: prod/db
        property: password
```

Пример безопасного service account и RBAC:

```yaml
# templates/rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "service.fullname" . }}

---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ include "service.fullname" . }}
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ include "service.fullname" . }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ include "service.fullname" . }}
subjects:
  - kind: ServiceAccount
    name: {{ include "service.fullname" . }}
```
