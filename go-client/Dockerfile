FROM golang:1.24-alpine AS dev
ENV TZ="Europe/Warsaw"
RUN apk update && apk upgrade
RUN apk add --no-cache bash git openssh autoconf automake libtool gettext gettext-dev make g++ texinfo curl
RUN mkdir /go/pkg
RUN chmod a+rwx /go/pkg
ARG APP_UID=1000
RUN adduser -s /bin/sh -u ${APP_UID} -D appuser
USER appuser
WORKDIR /go/src/go-client

FROM dev AS builder
USER root
COPY . /go/src/go-client
WORKDIR /go/src/go-client
RUN go build -buildvcs=false -o bin/go-client ./

FROM alpine:3.21.3 AS prod
ARG APP_UID=1000
RUN adduser -s /bin/sh -u ${APP_UID} -D appuser
USER appuser
COPY --from=builder /go/src/go-client/bin/go-client /usr/local/bin/go-client