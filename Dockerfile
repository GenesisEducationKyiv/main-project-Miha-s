FROM golang:1.20.4-alpine as build-stage

WORKDIR /go/src/btc-test-task

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /btc-test-task .

RUN apk --no-cache add ca-certificates

FROM scratch

COPY --from=build-stage /btc-test-task /btc-test-task
COPY --from=build-stage /go/src/btc-test-task/.env /.env

COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 2777

ENTRYPOINT ["/btc-test-task"]