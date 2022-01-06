FROM golang:1.17

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install vim gdb sudo dbus curl -y

# Setup dlv for debugging
RUN go get -d github.com/go-delve/delve/cmd/dlv

ARG DIR=/go/src/github.com/pranitbauva1997/splitwise-demo
WORKDIR ${DIR}
COPY go.mod go.sum ${DIR}
COPY . .

RUN make build

ENTRYPOINT ["./web"]
EXPOSE 8080