FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 5000

CMD ["go", "run", "cmd/main.go"]