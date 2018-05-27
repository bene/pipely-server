FROM golang

ADD . /go/src/github.com/bene/wetube-server
WORKDIR /go/src/github.com/bene/wetube-server

RUN go get
RUN go install

ENV HOST_DOMAIN=localhost

EXPOSE 8080

ENTRYPOINT /go/bin/wetube-server