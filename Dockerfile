FROM golang

ADD . /go/src/github.com/bene/wetube-server
WORKDIR /go/src/github.com/bene/wetube-server

RUN go get
RUN go install

ENV ADDRESS=:6550

EXPOSE 6550

ENTRYPOINT /go/bin/wetube-server