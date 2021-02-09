# Atlas

### Thingspect IoT platform. It's state-of-the-art.

## Getting Started

```
docker-compose -f build/docker-compose.yml up
make test
RACE=y make test
```

## Running an API Locally

First complete the above steps. Then:

```
atlas-create org testorg testadmin@thingspect.com testpass
API_PWT_KEY=$(dd if=/dev/random bs=1 count=32|base64) api

curl -v -X POST -d '{"email":"testadmin@thingspect.com", "orgName":"testorg", \
"password":"testpass"}' http://127.0.0.1:8000/v1/sessions/login
```

OpenAPI live docs are available at
[http://localhost:8000/](http://localhost:8000/).

## Use of Build Tags In Tests

All non-generated test files should have build tags, including `main_test.go`.
Due to limitations of developer tools and extensions, negated tags are used.

For example, to tag a file as a unit test:

```
// +build !integration
```

To tag a file as an integration test:

```
// +build !unit
```

To find test files that are missing build tags, the following command can be
run:

```
find . -name '*_test.go' -type f|grep -v /mock_|xargs grep -L '// +build'
```
