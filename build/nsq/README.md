# nsq

GitHub Actions
[do not support](https://github.community/t/how-do-i-properly-override-a-service-entrypoint/17435)
overriding a Docker `CMD` or `ENTRYPOINT` for services. These images wrap the
production NSQ images and set specific `CMD` values for each service.

## Build

```
docker login -u thingspect

docker build -f Dockerfile-nsqlookupd -t thingspect/nsqlookupd:v1.2.0 .
docker push thingspect/nsqlookupd:v1.2.0

docker build -f Dockerfile-nsqd -t thingspect/nsqd:v1.2.0 .
docker push thingspect/nsqd:v1.2.0

docker build -f Dockerfile-nsqadmin -t thingspect/nsqadmin:v1.2.0 .
docker push thingspect/nsqadmin:v1.2.0

docker logout
```

## Usage

```
docker run -it --env LOG_LEVEL=info thingspect/nsqlookupd:v1.2.0

docker run -it --env LOOKUP_ADDR=nsqlookupd:4160 --env LOG_LEVEL=info thingspect/nsqd:v1.2.0

docker run -it --env LOOKUP_ADDR=nsqlookupd:4161 --env LOG_LEVEL=info thingspect/nsqadmin:v1.2.0
```
