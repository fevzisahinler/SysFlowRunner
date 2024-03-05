FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o sysflowrunner main.go

FROM ubuntu:latest

COPY .env .

COPY --from=builder /app/sysflowrunner /usr/local/bin/sysflowrunner

CMD [ "sysflowrunner" ]