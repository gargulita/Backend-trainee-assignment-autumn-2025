FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o pr-reviewer ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/pr-reviewer .

EXPOSE 8080

CMD ["./pr-reviewer"]
