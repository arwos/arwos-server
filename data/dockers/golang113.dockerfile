FROM golang:1.13-alpine
ENV GO111MODULE=on
RUN apk update && \
    apk add --virtual build-dependencies build-base \
    bash git
WORKDIR /app