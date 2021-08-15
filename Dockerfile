FROM golang:1.16

RUN mkdir /app
WORKDIR /app

COPY ./api/go.mod ./api/go.sum /app/
RUN go mod tidy
RUN go mod download
