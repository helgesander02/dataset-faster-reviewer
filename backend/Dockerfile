# Builder Stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Final Stage - Use a glibc-based image instead of Alpine
FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /app/main .
COPY ./config ./config
EXPOSE 8080
CMD ["./main"]