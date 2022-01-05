FROM golang:1.17

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install vim gdb sudo dbus curl

# Setup dlv for debugging
RUN go get -d github.com/go-delve/delve/cmd/dlv

