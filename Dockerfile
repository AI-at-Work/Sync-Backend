FROM golang:1.22.5-alpine

WORKDIR /go/src/sync_worker

ADD . .

RUN go mod tidy

RUN go build -o main main.go

CMD ["./main"]
