# drone-deb WIP

Drone plugin for simple debian dpkg (.deb) file packaing. For usage information
and a listing of the available options please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following command:

```sh
go build
```

## Run tests

Run the tests locally using drone cli's exec command:

```sh
drone exec
```

## Docker

Build the docker image with the following commands:

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build -t plugins/deb .
```
