FROM golang:alpine

# git is required for fetching dependencies
RUN apk update && apk add --no-cache git

COPY . /go/src/github.com/paulcarmichael/fincrypt.io

RUN go get -d -v /go/src/github.com/paulcarmichael/fincrypt.io
RUN go install -v /go/src/github.com/paulcarmichael/fincrypt.io

ENTRYPOINT /go/bin/fincrypt.io

EXPOSE 80
