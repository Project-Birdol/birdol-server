FROM golang:latest

RUN mkdir /app
WORKDIR /app

RUN go mod init github.com/MISW/birdol-server 
RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-sql-driver/mysql
RUN go get gorm.io/gorm