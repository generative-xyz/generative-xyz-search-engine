FROM golang:1.19-alpine AS builder
WORKDIR /go/src/search-engine-worker

RUN apk update && apk upgrade && \
    apk --update add git gcc make libc-dev openssh linux-headers

COPY . .
RUN go mod tidy
RUN go build -o build/search-engine-worker cmd/worker/*.go

FROM alpine as release
WORKDIR /app

COPY --from=builder /go/src/search-engine-worker/build/search-engine-worker search-engine-worker
COPY --from=builder /go/src/search-engine-worker/env env

ENTRYPOINT ["./search-engine-worker"]
