FROM golang:1.20.4-alpine as build-stage

WORKDIR /go/src/logs-client

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /logs-client .

FROM alpine:3.17

COPY --from=build-stage /logs-client /logs-client
COPY --from=build-stage /go/src/logs-client/keep-logs.sh /keep-logs.sh

RUN apk add --no-cache bash
RUN chmod +x /keep-logs.sh

CMD /keep-logs.sh
