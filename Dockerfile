FROM golang:1.9-alpine

WORKDIR /go/src/github.com/ultrabluewolf/manifest-manager

ADD . /go/src/github.com/ultrabluewolf/manifest-manager/

RUN apk update
RUN apk add git

RUN go get -u github.com/golang/dep/cmd/dep

RUN go install ./cmd/manifest-manager
