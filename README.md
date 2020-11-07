# Thingspect Atlas

## Getting Started

```
docker-compose up
make test
```

Optional:

```
RACE=y make test
```

## Use of Build Tags In Tests

All test files should have build tags, including `main_test.go`. Due to
limitations of some developer tools and plugins, negated tags are used.

For example, to tag a file as a unit test:

```
// +build !integration
```

To tag as an integration test:

```
// +build !unit
```

To find test files that are missing build tags, the following command can be
run:

`find . -type f -name *_test.go|xargs grep -L '// +build'`
