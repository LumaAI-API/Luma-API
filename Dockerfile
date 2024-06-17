FROM golang AS builder

ENV GO111MODULE=on \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w -extldflags '-static'" -o lumaApi

FROM alpine:latest

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata gcc

COPY --from=builder /build/lumaApi /

EXPOSE 8000

ENTRYPOINT ["/lumaApi"]
