FROM golang:1.19-alpine

RUN apk update && apk upgrade
RUN apk add --no-cache git make protobuf protobuf-dev
RUN git clone --depth=1 https://github.com/googleapis/googleapis.git \
/tmp/googleapis
RUN git clone --depth=1 https://github.com/grpc-ecosystem/grpc-gateway.git \
/tmp/grpc-gateway
RUN git clone --depth=1 https://github.com/envoyproxy/protoc-gen-validate.git \
/tmp/pgv

WORKDIR /api
