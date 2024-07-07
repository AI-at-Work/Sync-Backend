FROM golang:1.22.5-alpine

# Install PostgreSQL client
RUN apk add --no-cache postgresql-client

WORKDIR /go/src/sync_worker

ADD . .

RUN chmod +x wait-for-it.sh

RUN go mod tidy

RUN go build -o main main.go
