# Atlas

### Thingspect IoT platform

## Getting Started

Install [Docker](https://docs.docker.com/get-started/overview/),
[Docker Compose](https://docs.docker.com/compose/), and
[Go](https://go.dev/dl/).

```
docker compose -f build/docker-compose.yml up -d
make test
RACE=y make test
docker compose -f build/docker-compose.yml down
```

## Running an API Locally

First complete the above steps. Then:

```
atlas-create org testorg testadmin@thingspect.com testpass
API_PWT_KEY=$(dd if=/dev/random bs=1 count=32|base64) API_API_HOST=127.0.0.1 atlas-api

curl -v -X POST -d '{"email":"testadmin@thingspect.com", "orgName":"testorg", "password":"testpass"}' http://localhost:8000/v1/sessions/login
```

OpenAPI live docs are available at
[http://localhost:8000/](http://localhost:8000/).

## Deploying

[Docker Compose](https://docs.docker.com/compose/) files for the Atlas platform
and its dependencies are available in `build/deploy/`. These can be used for a
single-system deploy, or as templates for orchestration tooling such as
[Nomad](https://www.nomadproject.io/) or [Kubernetes](https://kubernetes.io/).
Keys should be provided where applicable.

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
find . -name '*_test.go' -type f|grep -v /mock_|xargs grep -L '//go:build'
```
