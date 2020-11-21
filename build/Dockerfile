# docker build -f build/Dockerfile -t ghcr.io/thingspect/atlas:test .
FROM alpine

# Add CA certificates and timezone info
RUN apk update && apk upgrade
RUN apk add --no-cache ca-certificates tzdata

# Add build files
COPY bin/mqtt-ingestor* /bin/

USER nobody
