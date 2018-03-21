FROM golang:1.10-alpine as builder

ARG VERSION

RUN apk add --no-cache git ca-certificates

RUN git clone --branch "v1.0" --single-branch --depth 1 \
    https://github.com/korylprince/fileenv.git /go/src/github.com/korylprince/fileenv

RUN git clone --branch "$VERSION" --single-branch --depth 1 \
    https://github.com/korylprince/bisd-device-checkin-server.git  /go/src/github.com/korylprince/bisd-device-checkin-server

RUN go install github.com/korylprince/fileenv
RUN go install github.com/korylprince/bisd-device-checkin-server

FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/fileenv /
COPY --from=builder /go/bin/bisd-device-checkin-server /
COPY setenv.sh /

CMD ["/fileenv", "sh", "/setenv.sh", "/bisd-device-checkin-server"]
