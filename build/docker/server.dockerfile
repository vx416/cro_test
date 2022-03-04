
FROM golang:1.16 AS builder
WORKDIR /cro_test
ENV GO111MODULE=on 

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -o main

FROM alpine:3.15
ARG BUILD_TIME
ARG SHA1_VER

RUN apk update && \
    apk upgrade && \
    apk add --no-cache curl tzdata && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /cro_test
COPY --from=builder /cro_test/main /cro_test/main
COPY ./configs/dev.yaml ./configs/prod.yaml

RUN ls
ENV SHA1_VER=${SHA1_VER}
ENV BUILD_TIME=${BUILD_TIME}
ENV CONFIG_NAME=prod
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /cro_test
USER appuser

CMD ["./main", "server"]
