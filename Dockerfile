FROM golang:1.19-alpine as builder 

RUN apk update \
  && apk add --no-cache git curl make gcc g++ tzdata

WORKDIR /app

COPY ./api .

RUN go mod tidy
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -tags netgo -o ./birdol-server

FROM gcr.io/distroless/static-debian11 as production

COPY --from=builder /app/birdol-server /

ENV PORT=${PORT}
CMD ["/birdol-server"]
