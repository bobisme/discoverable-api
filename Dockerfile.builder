FROM golang:1.11
WORKDIR /app
VOLUME /app
RUN go get github.com/Jun-Chang/gdbinder
RUN go get github.com/tinylib/msgp
RUN go get github.com/swaggo/swag/cmd/swag
