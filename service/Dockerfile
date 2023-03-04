FROM golang:1.17-alpine
WORKDIR /app

COPY server.go .

RUN go build server.go

EXPOSE 8080 8081 8082
CMD ["./server"]
