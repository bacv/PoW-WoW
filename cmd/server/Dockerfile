# syntax=docker/dockerfile:1
FROM golang:1.18-bullseye
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /server /app/cmd/server
EXPOSE 8080
CMD ["/server"]

