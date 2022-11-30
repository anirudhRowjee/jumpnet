# syntax=docker/dockerfile:1

FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN go build -o /kafka-broker

# run Command to come in docker-compose file
