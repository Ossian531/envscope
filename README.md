# envscope

Two tiny Go services:

- **backend/** — on `GET /`, returns its process environment variables as JSON.
- **frontend/** — fetches the backend and renders the env vars on a styled web page (with live filtering).

## Run

In one terminal:

```sh
cd backend
SOME_VAR=hello ANOTHER=world go run .
# backend listening on :8080
```

In another:

```sh
cd frontend
go run .
# frontend listening on :3000 (backend: http://localhost:8080/)
```

Open http://localhost:3000.

## Configuration

| Service  | Env var         | Default                   |
|----------|-----------------|---------------------------|
| backend  | `BACKEND_ADDR`  | `:8080`                   |
| frontend | `FRONTEND_ADDR` | `:3000`                   |
| frontend | `BACKEND_URL`   | `http://localhost:8080/`  |

## API

`GET /` on the backend:

```json
{
  "hostname": "myhost",
  "served_at": "2026-06-18T12:00:00Z",
  "count": 42,
  "env_vars": { "PATH": "/usr/bin", "SOME_VAR": "hello" }
}
```

> Note: this intentionally exposes the backend's full environment. Don't point a
> public frontend at a backend holding secrets.
# envscope
