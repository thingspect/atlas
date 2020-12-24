# Atlas

### Thingspect IoT platform. It's state-of-the-art.

## Getting Started

```
docker-compose -f build/docker-compose.yml up
make test
```

Optional:

```
RACE=y make test
```

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

`find . -type f -name \*_test.go|grep -v /mock_|xargs grep -L '// +build'`
