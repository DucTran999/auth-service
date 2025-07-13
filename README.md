# 🛡️ Auth Service

[![Go Report Card](https://goreportcard.com/badge/github.com/DucTran999/auth-service)](https://goreportcard.com/report/github.com/DucTran999/auth-service)
[![Go](https://img.shields.io/badge/Go-1.23-blue?logo=go)](https://golang.org)
[![codecov](https://codecov.io/gh/DucTran999/auth-service/branch/master/graph/badge.svg)](https://codecov.io/gh/DucTran999/auth-service)
[![Known Vulnerabilities](https://snyk.io/test/github/ductran999/auth-service/badge.svg)](https://snyk.io/test/github/ductran999/auth-service)
[![License](https://img.shields.io/github/license/DucTran999/auth-service)](LICENSE)

A personal project exploring secure session-based authentication in Go, using Redis, PostgreSQL, and Clean Architecture.

> ✅ Designed for learning.  
> ⚙️ Built like a real-world system.

---

## 📘 About This Project

This project was created to deepen my understanding of authentication system design, session lifecycle management, and Go service architecture. It is not intended for production, but it reflects production-like patterns.

---

## 🚀 Features

- 🔐 Login / Logout support with both session cookies and JWT
- 🍪 Session-Based Auth using Redis (fast cache) + PostgreSQL (durable storage)
- 🔑 JWT-Based Auth with access/refresh tokens and RSA key signing
- 🧠 Clean Architecture: Handler → UseCase → Repository separation
- 📦 Microservice-Friendly: Shared session store & stateless JWT integration
- ⏱️ Configurable session TTL, HTTP-only cookies, and IP/User-Agent tracking
- 📜 Full audit trail via PostgreSQL for login and session activity

## Project Structure


```bash
auth-service/
├── cmd/                          # Application entry point (main.go)
│   └── main.go                   # Starts the service with DI container
│
├── config/                       # Application configuration loading using Viper
│   ├── config.go
│   └── loader.go
│
├── gen/                          # OpenAPI generated types and server interface
│   └── http/
│       ├── server.gen.go
│       └── types.gen.go
│
├── environment/                  # Development environment setup (migrations, compose)
│   ├── migration/                # Database migration files
│   └── docker-compose.yaml       # Compose file to start PostgreSQL and Redis
│
├── internal/                     # Internal app modules (not exposed externally)
│   │
│   ├── container/                # DI container: setup DB, Redis, repos, usecases
│   │
│   ├── server/                   # HTTP server, router setup, validator registration
│   │   ├── http_server.go
│   │   ├── router.go
│   │   └── validator.go
│   │
│   ├── worker/                   # Background tasks (e.g., cleanup jobs)
│   │   └── session_cleaner.go
│   │
│   ├── handler/                  # HTTP handlers (controllers)
│   │   ├── rest/                 # Session-based + JWT handler logic
│   │   └── auth_handler.go
│   │
│   ├── usecase/                  # Business logic (interactors)
│   │   ├── port/                 # Interfaces to handler & repository
│   │   ├── dto/                  # Data transfer objects
│   │   └── auth_usecase.go       # Auth logic implementation
│   │
│   ├── domain/                   # Entities, enums, and domain-level interfaces
│   │   └── account.go
│   │
│   └── repository/               # Data persistence logic
│       ├── account_repo.go
│       └── session_repo.go
│
├── scripts/                      # Scripts for automation and testing
│   ├── testenv.sh                # Setup test environment with Redis/Postgres
│   ├── unittest.sh               # Run unit tests
│   ├── integration.sh            # Run integration tests (real DB)
│   ├── api-test.sh               # Run API tests (end-to-end flow)
│   └── gen-keys.sh               # Generate RSA key pairs for JWT
│
├── pkg/                          # Shared packages (reusable across layers)
│   ├── cache.go                  # Redis caching utilities
│   ├── hasher.go                 # Password hashing using Argon2
│   └── signer.go                 # JWT signing and parsing
│
├── go.mod                        # Go module file
├── go.sum                        # Module checksums
├── Dockerfile                    # Docker build definition
└── README.md                     # Project documentation
```

---

## 📚 API Endpoints

### 🔐 Auth Endpoints

#### 🧾 Session-Based Auth (session_id cookie)

Sessions are stored as secure, HTTP-only cookies.

| Method | Endpoint                   | Description                     | Auth Required           |
| ------ | -------------------------- | ------------------------------- | ----------------------- |
| POST   | `/api/v1/register`         | Register a new user account     | ❌ No                   |
| POST   | `/api/v1/login`            | Authenticate and create session | ❌ No                   |
| DELETE | `/api/v1/logout`           | Destroy the session and logout  | ❌ No _(uses cookie)_   |
| PATCH  | `/api/v1/account/password` | Change account password         | ✅ Yes (session cookie) |

#### 🔑 JWT-Based Auth (access_token + refresh_token)

JWT-based flow using Authorization: Bearer <access_token> and refresh_token via secure cookie.

| Method | Endpoint                | Description                               | Auth Required              |
| ------ | ----------------------- | ----------------------------------------- | -------------------------- |
| POST   | `/api/v2/login`         | Login and receive access + refresh tokens | ❌ No                      |
| POST   | `/api/v2/token/refresh` | Refresh tokens                            | ❌ No _(relies on cookie)_ |
| POST   | `/api/v2/logout`        | Logout and revoke refresh token           | ❌ No _(relies on cookie)_ |

---

### ⚙️ Infrastructure Endpoints

| Method | Endpoint | Description                   | Auth Required |
| ------ | -------- | ----------------------------- | ------------- |
| GET    | `/livez` | Liveness probe (health check) | ❌ No         |

---

## ✅ Prerequisites

Before running this project, make sure you have the following tools installed:

| Tool                                  | Purpose                      | Install Command / Link                                   |
| ------------------------------------- | ---------------------------- | -------------------------------------------------------- |
| [Go](https://golang.org/dl/) (v1.20+) | Build and run the service    | [Download manually](https://go.dev/dl)                   |
| [Docker](https://www.docker.com/)     | Run Redis/PostgreSQL locally | [Download manually](https://docs.docker.com/get-docker/) |
| [Task](https://taskfile.dev)          | Simplified task runner       | `go install github.com/go-task/task/v3/cmd/task@latest`  |
| Make (optional)                       | Fallback task runner (Linux) | `sudo apt install make` or `brew install make`           |

> 📝 Note: This project uses [`task`](https://taskfile.dev) to automate common commands like running the service and setting up the test environment.

## 🧪 Running Locally

1. Copy the environment file and fill in required values:

```bash
cp .env.example .env
```

2. Start Redis, PostgreSQL, and run database migrations:

```bash
task testenv
```

3. Generate RSA keys:

```bash
task keys
```

4. Run the auth service:

```bash
task run
```

---

## 🧪 Testing

This project uses Go's built-in testing framework with mocks and table-driven tests.

### ✅ Run All Unit Tests

```bash
task unit_test
```

### ✅ Run All Integration Tests

- setup db and redis for testenv

```bash
task testenv
```

- run scripts integration test

```bash
task integration_test
```

### ✅ Run All API Tests

- Note: _if you already run it in integration tests, you can skip this step._

```bash
task testenv
```

- run scripts api test

```bash
task api_test
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

## 🤝 Contributing

We welcome PRs for:

- New features (e.g. 2FA, OAuth, etc.)
- Bug fixes
- Documentation improvements
- Test coverage

Please follow conventional commits and open an issue first for major changes.
