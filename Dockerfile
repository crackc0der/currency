FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o main cmd/currency/main.go cmd/currency/init.go

CMD ["/app/main"]
