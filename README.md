## PDF Payroll POC (Go)

Minimal Go HTTP service that fetches mock payroll JSON and generates a PDF pay stub via `/generate`.

### Prerequisites
- Go 1.20+

### Install deps
```bash
cd /Users/agionfriddo/pdf-poc
go mod tidy
```

### Run
```bash
PORT=9090 go run .
```

### Endpoints
- `GET /healthz` — health check
- `GET /generate` — returns a generated `application/pdf` attachment

### Example
```bash
curl -sS -D - http://localhost:9090/generate -o paystub.pdf
open paystub.pdf
```

### Structure
- `internal/models` — data models
- `internal/data` — mock JSON provider simulating API fetch
- `internal/pdf` — gofpdf generator
- `main.go` — HTTP server


# go-pdf-generator
