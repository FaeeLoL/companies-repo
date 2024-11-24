# Stage 1: Builder
FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o companies-repo ./cmd

# Stage 2: Production
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/companies-repo .

EXPOSE 8080

CMD ["sh", "-c", "${MIGRATE_BEFORE_START:-false} && ./companies-repo migrate-db --config configs/config.yaml; ./companies-repo http --config configs/config.yaml"]

