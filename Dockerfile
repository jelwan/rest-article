FROM golang:1.14 as builder
ARG PACKAGE
ADD . /go/src/$PACKAGE
WORKDIR /go/src/$PACKAGE
RUN APP_NAME=app CGO_ENABLED=0 GOOS=linux TARGET_FLAGS="-a -installsuffix cgo" make compile && chmod +x /go/src/$PACKAGE/bin/app

FROM alpine:latest
ARG PACKAGE
COPY --from=builder /go/src/$PACKAGE/bin/app /app
COPY --from=builder /go/src/$PACKAGE/data/config/app.yaml /data/config/app.yaml
EXPOSE 8080
CMD ["/app"]
