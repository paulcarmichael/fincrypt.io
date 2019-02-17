FROM golang:alpine

# git is required for fetching dependencies
RUN apk update && apk add --no-cache git

COPY . /go/src/github.com/paulcarmichael/fincrypt_http

RUN go get -d -v /go/src/github.com/paulcarmichael/fincrypt_http
RUN go install -v /go/src/github.com/paulcarmichael/fincrypt_http

ENTRYPOINT /go/bin/fincrypt_http

EXPOSE 80
