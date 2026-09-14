FROM golang:1-alpine

ARG PROTOC_VERSION=36.1

RUN apk update && apk upgrade
RUN apk add --no-cache git make curl unzip
RUN curl -L -o /tmp/protoc.zip \
https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-aarch_64.zip
RUN unzip /tmp/protoc.zip -d /usr
# Required for imports
RUN git clone --depth=1 https://github.com/googleapis/googleapis.git \
/tmp/googleapis
RUN git clone --depth=1 https://github.com/grpc-ecosystem/grpc-gateway.git \
/tmp/grpc-gateway
RUN git clone --depth=1 https://github.com/bufbuild/protovalidate.git \
/tmp/pv

WORKDIR /proto
