FROM golang:alpine

COPY . /go/src/github.com/paulcarmichael/fincrypt_http

RUN go install -v /go/src/github.com/paulcarmichael/fincrypt_http

ENTRYPOINT /go/bin/fincrypt_http

EXPOSE 80