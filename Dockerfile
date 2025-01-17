FROM golang:1.22-alpine

WORKDIR /

COPY go.mod ./

RUN go mod download

COPY . .
RUN go build -o bin/server main.go
EXPOSE 3000

CMD ["./bin/server"]
