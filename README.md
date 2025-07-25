# 🧊 freezeRepo
---

## 🚀 Быстрый старт

### 1. Клонируй репозиторий

```bash
git clone -b master https://github.com/ttiiuus/freezeRepo.git
cd freezeRepo
```
###2. Установи зависимости
```
go mod tidy
```
###3. Создай файл config.yaml с содержимым:
```
host: "localhost"
port: 8080

database:
  host: "localhost"
  port: 5432
  user: "your_user"
  password: "your_password"
  dbname: "your_db"
  sslmode: "disable"
  mongo:
    uri: "mongodb://localhost:27017"
    name: "your_mongo_db"

auth:
  jwt_secret: "your_secret_key"
  #Заменить значения your_user, your_password, your_db, your_mongo_db и your_secret_key на свои реальные настройки.
  ```


###4. Запусти сервер
```
go run ./cmd/server
```
