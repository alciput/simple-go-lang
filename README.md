# Go Auth API (SQLite + JWT)

Production-ready backend: register, login, and protected profile. Uses **Gin**, **SQLite**, **JWT**, **Swagger**, and a consistent API response envelope. No Docker required.

## Run locally

1. **Install deps and build**

   Using the Makefile (recommended if your project lives inside GOPATH, e.g. `~/Go` on macOS):

   ```bash
   make tidy
   make build
   ```

   Or run Go directly (if you get "go.mod file not found", use the Makefile or see [Troubleshooting](#troubleshooting) below):

   ```bash
   go mod tidy
   go build -o app .
   ```

2. **Run** (defaults: port 8080, DB file `app.db`, dev JWT secret)

   ```bash
   ./app
   ```

   Or with env:

   ```bash
   PORT=8080 JWT_SECRET=your-secret DB_PATH=app.db ./app
   ```

3. **Swagger UI:** `http://localhost:8080/swagger/index.html`

4. **Tests**

   ```bash
   make test
   ```
   or `go test ./...`

## Troubleshooting

**"go.mod file not found" or "ignoring go.mod in $GOPATH"**  
If your project is inside GOPATH (e.g. `/Users/you/Go` on macOS, where `Go` and `go` are the same folder), Go may ignore the module. Fix:

- Use the Makefile: `make tidy` then `make build`, or  
- Run with GOPATH unset: `GOPATH= go mod tidy` and `GOPATH= go build -o app .`

To fix it permanently, set GOPATH to a path that does not contain this project, e.g. add to `~/.zshrc`:  
`export GOPATH=$HOME/gopath`

## API response format

All responses use an envelope:

- Success: `{ "data": <payload>, "error": "" }`
- Error: `{ "data": null, "error": "<message>" }`

## API

Base URL: `http://localhost:8080` (or your `PORT`).

### Register

```bash
curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Example response (201):**

```json
{ "data": { "id": 1, "email": "user@example.com" }, "error": "" }
```

### Login

```bash
curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Example response (200):**

```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": { "id": 1, "email": "user@example.com", "created_at": "2025-01-31T12:00:00Z" }
  },
  "error": ""
}
```

### Profile (protected)

Use the `token` from login in the `Authorization` header:

```bash
TOKEN="<paste token from login>"
curl -s http://localhost:8080/profile -H "Authorization: Bearer $TOKEN"
```

**Example response (200):**

```json
{
  "data": { "id": 1, "email": "user@example.com", "created_at": "2025-01-31T12:00:00Z" },
  "error": ""
}
```

**Without token or invalid token (401):**

```json
{ "data": null, "error": "missing authorization header" }
```

## Security notes (see comments in code)

- **Passwords:** bcrypt only; never stored or logged in plain text.
- **JWT:** Signed with `JWT_SECRET`; 24h expiry. Set a strong secret in production.
- **Errors:** Generic messages only (e.g. "invalid email or password") to avoid user enumeration.
- **SQL:** All queries use prepared statements; no string concatenation for user input.
- **Config:** Port, JWT secret, and `DB_PATH` from env; no secrets in code.

## Project layout

```
.
├── main.go             # entrypoint, route wiring, Swagger
├── config/             # env-based config
├── docs/               # Swagger spec (docs.go); regenerate with: swag init -g main.go --parseDependency --parseInternal
├── handlers/           # register, login, profile (+ _test.go)
├── middleware/         # JWT auth
├── models/             # User
├── response/           # envelope + concrete response structs
├── store/              # UserRepository interface
├── store/sqlite/       # SQLite implementation
├── store/mock.go       # in-memory mock for tests
└── go.mod
```
