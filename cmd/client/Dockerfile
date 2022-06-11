# syntax=docker/dockerfile:1
FROM golang:1.18-bullseye
WORKDIR /app
COPY . . 
RUN go mod download
RUN go build -o /client /app/cmd/client 
CMD ["/client"]

