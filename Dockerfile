FROM golang

ADD . /go/src/github.com/bene/wetube-server
WORKDIR /go/src/github.com/bene/wetube-server

RUN go get -u github.com/golang/dep/...
RUN dep ensure
RUN go install

ENV ADDRESS=:6550

EXPOSE 6550

ENTRYPOINT /go/bin/wetube-server