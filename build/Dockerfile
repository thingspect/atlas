# Manual build:
# docker login -u thingspect -p XXX ghcr.io
# TAG=$(git rev-parse --short=8 HEAD)
# docker buildx create --use
## docker buildx build -f build/Dockerfile -t ghcr.io/thingspect/atlas:${TAG} --platform linux/amd64,linux/arm64 --push .
# docker buildx build -f build/Dockerfile -t ghcr.io/thingspect/atlas:${TAG} --push .
# docker buildx rm
# docker logout ghcr.io
FROM golang:1.18-alpine AS build

# Add packages needed to support builds and tests
RUN apk update && apk upgrade
RUN apk add --no-cache make git
RUN go version
RUN go env

COPY . /opt/atlas/
RUN mkdir -p /opt/atlas/bin
WORKDIR /opt/atlas/

RUN GOBIN=/opt/atlas/bin make install

FROM alpine

# Add CA certificates and timezone info
RUN apk update && apk upgrade
RUN apk add --no-cache ca-certificates tzdata

# Add build files
COPY --from=build /opt/atlas/bin/* /bin/
COPY web/* /web/

USER nobody
