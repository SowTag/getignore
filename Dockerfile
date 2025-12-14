# Builder stage
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o getignore ./cmd/getignore/main.go

# Runtime stage
FROM alpine:latest

LABEL authors="Maddock"

# HTTPS requests fix
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/getignore /usr/local/bin/getignore

ENTRYPOINT ["getignore"]