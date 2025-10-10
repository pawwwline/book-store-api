FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o book-store-api ./cmd/book-store-api/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/book-store-api .
COPY migrations ./migrations
COPY docs ./docs
COPY .env .

EXPOSE 8080

CMD ["./book-store-api"]
