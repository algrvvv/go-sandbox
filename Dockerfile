FROM golang:1.22-alpine

RUN apk add --no-cache docker-cli

WORKDIR /app

COPY . .

RUN mkdir -p /tmp/go-sandbox/

RUN go mod download

RUN go build -o sandbox .

RUN go install golang.org/x/tools/cmd/goimports@latest

COPY go-runner.tar /go-runner.tar
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]

CMD ["./sandbox"]
