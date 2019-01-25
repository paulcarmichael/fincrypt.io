FROM golang:alpine

COPY . /go/src/github.com/paulcarmichael/fincrypt

RUN go install -v /go/src/github.com/paulcarmichael/fincrypt

ENTRYPOINT /go/bin/fincrypt

EXPOSE 80