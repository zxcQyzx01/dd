FROM golang:1.22-alpine

RUN apk add --no-cache curl netcat-openbsd

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o main ./user

EXPOSE 50053

CMD ["./main"] 