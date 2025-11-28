FROM golang:1.22-alpine AS builder

WORKDIR /src

# Install git for fetching modules that require it.
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/users-service ./cmd/server

FROM alpine:3.19

WORKDIR /app

RUN adduser -D appuser
USER appuser

COPY --from=builder /bin/users-service /app/users-service

EXPOSE 8080

CMD ["/app/users-service"]
