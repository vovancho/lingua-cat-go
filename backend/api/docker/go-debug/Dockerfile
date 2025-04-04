FROM golang:1.24-alpine3.21

RUN go install github.com/go-delve/delve/cmd/dlv@v1.24
RUN go install github.com/air-verse/air@v1.61

WORKDIR /app

CMD ["air"]
