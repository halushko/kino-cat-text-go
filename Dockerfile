FROM golang:1.23.2 AS builder
WORKDIR /app
RUN go mod init kino-cat-text-go
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /app/kino-cat-text-go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/kino-cat-text-go .
#RUN chmod +x ./kino-cat-text-go
CMD ["./kino-cat-text-go"]
