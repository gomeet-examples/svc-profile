# svc-profile

Gomeetexamples's svc-profile service

This service is generated with [gomeet](https://github.com/gomeet/gomeet). The initial command is :

```shell
$ gomeet new github.com/gomeet-examples/svc-profile \
  --default-prefixes=svc-,gomeet-svc- \
  --proto-name=pb
```

## Up and running

- [docker usage](docs/docker/README.md)
- [docker-compose usage](docs/docker-compose/README.md)
- [baremetal](docs/baremetal/README.md)

## Usage

- [See svc-profile usage](docs/usage/README.md)
- [gRPC services](docs/grpc-services/README.md)

## Common development tasks

- [Working with the sources](docs/devel/working_with_the_sources/README.md)
- [Sub service declaration](docs/devel/add_sub_service/README.md)

## ROADMAP

- [x] Serve the grpc service
- [x] Serve the http service
- [x] Serve the grpc-gateway service
- [x] Multiplex grpc/http
- [x] Prometheus metrics
- [x] Add vendoring tool
  - [dep](https://github.com/golang/dep) win qualification vs glide and govendor
  - [retool](https://github.com/twitchtv/retool) for tools dependencies
- [x] Add serve cmd
- [x] Add cli cmd
- [x] Add console cmd
- [x] Generate/Embbed swagger.json
- [x] Add TLS support
- [x] Add docker environement (docker file, docker-compose)
- [x] Add prometheus to docker environement (docker-compose)
- [ ] Add Kubernetes
- [ ] Add a data store
- [x] Add unit tests
- [x] Add functional testing
- [x] Add load testing
- [ ] Add some clients examples (multilang: PHP, RUBY, Node.js, Python ...)
- [x] Add git workflow ([git flow](https://danielkummer.github.io/git-flow-cheatsheet/))
