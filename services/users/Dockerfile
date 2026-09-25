FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./
COPY pkg/xerrors/ ./pkg/xerrors/
COPY services/users/ ./services/users/

WORKDIR /app/services/users
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/user-service ./cmd/api

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /

COPY --from=builder /app/user-service /user-service

EXPOSE 50051
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/user-service"]
