FROM golang:1.19.3-alpine

USER nobody:nogroup

ARG GOOS

ENV CGO_ENABLED=0 XDG_CACHE_HOME=/tmp/.cache

WORKDIR /go/src/go-mailin8
COPY . .
RUN go install -v
