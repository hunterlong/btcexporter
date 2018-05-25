FROM golang

ADD . /go/src/github.com/hunterlong/btcexporter
RUN cd /go/src/github.com/hunterlong/btcexporter && go get
RUN go install github.com/hunterlong/btcexporter

ENV PORT 9019

RUN mkdir /app
WORKDIR /app
ADD addresses.txt /app

EXPOSE 9019

ENTRYPOINT /go/bin/btcexporter
