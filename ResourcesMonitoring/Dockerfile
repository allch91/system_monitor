FROM golang:alpine

RUN apk add --no-cache bash git
RUN go get github.com/thedevsaddam/renderer

COPY common /usr/local/go/src/common
COPY github.com /usr/local/go/src/github.com
COPY golang.org /usr/local/go/src/golang.org
COPY helpers /usr/local/go/src/helpers

RUN mkdir /app  
WORKDIR /app

COPY . .

RUN go build -o main .
EXPOSE 3001
CMD ["./main"]