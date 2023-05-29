FROM alpine:3.7 as nightwather
WORKDIR /
COPY ./cmd/bin/nightwatcher .
ENTRYPOINT ["./nightwatcher"]