FROM golang:1.4.2-cross

RUN go get github.com/mitchellh/gox
RUN go get github.com/aktau/github-release
RUN go get github.com/tools/godep
RUN go get golang.org/x/tools/cmd/cover

ENV GOPATH /go/src/github.com/docker/libcompose/Godeps/_workspace:/go
ENV COMPOSE_BINARY /go/src/github.com/docker/libcompose/docker-compose
ENV USER root

WORKDIR /go/src/github.com/docker/libcompose

ADD . /go/src/github.com/docker/libcompose
