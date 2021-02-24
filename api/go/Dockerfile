FROM golang:1.16-alpine

RUN apk update && apk upgrade
RUN apk add --no-cache make protobuf protobuf-dev git
RUN git clone --depth=1 https://github.com/googleapis/googleapis.git \
/tmp/googleapis
RUN git clone --depth=1 https://github.com/envoyproxy/protoc-gen-validate.git \
/tmp/pgv

WORKDIR /api
