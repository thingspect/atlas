# nsq

GitHub Actions
[do not support](https://github.community/t/how-do-i-properly-override-a-service-entrypoint/17435)
overriding a Docker `CMD` or `ENTRYPOINT` for services. These images wrap the
production NSQ images and set specific `CMD` values for each service.

## Build

```
docker login -u thingspect -p XXX ghcr.io
docker buildx create --use

docker buildx build -f Dockerfile-nsqlookupd -t ghcr.io/thingspect/nsqlookupd:v1.2.1 --platform linux/amd64,linux/arm64 --push .

docker buildx build -f Dockerfile-nsqd -t ghcr.io/thingspect/nsqd:v1.2.1 --platform linux/amd64,linux/arm64 --push .

docker buildx rm
docker logout ghcr.io
```

## Usage

```
docker run -it --env LOG_LEVEL=info ghcr.io/thingspect/nsqlookupd:v1.2.1

docker run -it --env LOOKUP_ADDR=nsqlookupd:4160 --env LOG_LEVEL=info ghcr.io/thingspect/nsqd:v1.2.1
```
