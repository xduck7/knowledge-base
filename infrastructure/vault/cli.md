### Простой пример: записать и прочитать секрет

Рассмотрим базовый KV‑секрет (ключ‑значение), например `database/password`.

1. Включаем KV‑секретный движок (в dev он обычно уже есть, но покажем явно):  
   ```bash
   vault secrets enable -path=secret kv
   ```

2. Записываем секрет:  
   ```bash
   vault kv put secret/app/db password="super-secret-password" user="app_user"
   ```

3. Читаем секрет CLI:  
   ```bash
   vault kv get secret/app/db
   ```
   В выводе будут поля `password` и `user`.

4. Чтение через HTTP‑API (например, из бэкенда):  
   - HTTP запрос:  
     ```http
     GET /v1/secret/data/app/db
     X-Vault-Token: <VAULT_TOKEN>
     ```
   - Ответ (KV v2):  
     ```json
     {
       "request_id": "…",
       "data": {
         "data": {
           "password": "super-secret-password",
           "user": "app_user"
         },
         "metadata": {
           "version": 1,
           "created_time": "…"
         }
       }
     }
     ```

Бэкенд‑сервис обычно:  
- аутентифицируется в Vault (по JWT, Kubernetes auth, AppRole и т.п.);
- получает токен с ограниченными правами;  
- по этому токену читает нужные пути (`secret/data/app/db`) и забирает поля.
