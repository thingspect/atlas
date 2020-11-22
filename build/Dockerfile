# Manual build (-race not supported when cross compiling):
# DOCKER=y GOOS=linux make
# docker login -u thingspect -p XXX ghcr.io
# TAG=$(git rev-parse --short=8 HEAD)
# docker build -f build/Dockerfile -t ghcr.io/thingspect/atlas:${TAG} .
# docker push ghcr.io/thingspect/atlas:${TAG}
# docker logout ghcr.io
FROM alpine

# Add CA certificates and timezone info
RUN apk update && apk upgrade
RUN apk add --no-cache ca-certificates tzdata

# Add build files
COPY bin/mqtt-ingestor* /bin/

USER nobody
