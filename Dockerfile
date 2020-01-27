FROM golang:latest

LABEL maintainer="Vass Mark <vmark0818@gmail.com>"

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]