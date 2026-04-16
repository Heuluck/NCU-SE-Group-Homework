# GopherTodo Backend

Go HTTP backend for the frontend API in `fe/openapi.yaml`.

## Run

```bash
cd be
go run ./cmd/server
```

The server listens on `http://localhost:8080` by default and stores tasks in
`be/data/tasks.json`.

Configuration:

- `PORT`: HTTP port, defaults to `8080`
- `DATA_FILE`: task storage path, defaults to `data/tasks.json`

## Frontend Integration

Start the frontend with the API base URL pointing at this server:

```bash
cd fe
$env:VITE_API_BASE_URL="http://localhost:8080"
pnpm dev
```

## API

- `GET /tasks`: list all tasks
- `POST /tasks`: create a task with JSON body `{"content":"..."}`
- `GET /tasks/{id}`: get one task
- `POST /tasks/{id}/complete`: mark a task as completed
- `DELETE /tasks/{id}`: delete a task

Task JSON shape:

```json
{
  "id": 1,
  "content": "write backend",
  "status": "pending",
  "created_at": "2026-04-16T08:00:00Z",
  "completed_at": null
}
```

## Test

```bash
cd be
go test ./...
```
