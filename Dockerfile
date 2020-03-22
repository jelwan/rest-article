FROM golang:1.14 as builder
ARG PACKAGE
ADD . /go/src/$PACKAGE
WORKDIR /go/src/$PACKAGE
RUN APP_NAME=app CGO_ENABLED=0 GOOS=linux TARGET_FLAGS="-a -installsuffix cgo" make compile && chmod +x /go/src/$PACKAGE/bin/app

WORKDIR $GOPATH/src/github.com/johan-lejdung/go-microservice-api-guide/rest-api
COPY ./ .
RUN GOOS=linux GOARCH=386 go build -ldflags="-w -s" -v
RUN cp rest-api /

FROM alpine:latest
COPY --from=builder /rest-api /
CMD ["/rest-api"]