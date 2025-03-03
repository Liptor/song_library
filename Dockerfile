FROM golang:1.22.4 AS builder

WORKDIR /app

COPY main.go go.mod go.sum ./

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o server

FROM alpine:latest

COPY --from=builder /app/server ./

COPY static/ ./static/

EXPOSE 80

CMD ["./server", "--port", "80"]