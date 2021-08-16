FROM golang:1.16-alpine as builder 

RUN apk update \
  && apk add --no-cache git curl make gcc g++ 

WORKDIR /app

COPY ./api .

RUN go mod tidy
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /main

FROM alpine:3.9

COPY --from=builder /main .

ENV PORT=${PORT}
ENTRYPOINT ["/main"]
