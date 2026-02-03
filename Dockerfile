

FROM golang:1.25.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# Build with CGO disabled for Alpine compatibility
RUN CGO_ENABLED=0 go build -o /app/read-api ./cmd/read-api/main.go
RUN CGO_ENABLED=0 go build -o /app/write-api ./cmd/write-api/main.go
RUN CGO_ENABLED=0 go build -o /app/worker ./cmd/worker/main.go

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/read-api /app/read-api
COPY --from=builder /app/write-api /app/write-api
COPY --from=builder /app/worker /app/worker
RUN chmod +x /app/read-api /app/write-api /app/worker