To‑Do List (Go + React + TS)

Как запустить

1. **Запусти PostgreSQL (Docker):**
   ```sh
   docker run --name wails-todo-postgres -e POSTGRES_PASSWORD=todo123 -e POSTGRES_USER=todo -e POSTGRES_DB=wailstodo -p 5432:5432 -d postgres:16
   ```
2. **Создай таблицу задач:**
   ```sh
   psql -h localhost -U todo -d wailstodo -f internal/repo/migrations.sql
   ```
3. **Экспортируй переменную окружения:**
   ```sh
   export POSTGRES_DSN="postgres://todo:todo123@localhost:5432/wailstodo?sslmode=disable"
   ```
4. **Собери и запусти backend:**
   ```sh
   go build
   ./wails-todo
   ```
5. **Запусти frontend:**
   ```sh
   cd frontend
   npm install
   npm run dev
   ```
   Открой http://localhost:5173


Что реализовано

- [x] Добавление, удаление, отметка задач
- [x] Сохранение задач в PostgreSQL
- [x] repo → service → usecase (чистая архитектура)
- [x] Фильтрация и сортировка
- [x] Всё работает после перезапуска

