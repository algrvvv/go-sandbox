FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/runner runner.go

ENTRYPOINT ["runner"]
