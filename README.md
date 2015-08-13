# libcompose

A Go library for Docker Compose. It does everything the command-line tool does, but from within Go -- read Compose files, start them, scale them, etc.

**Note: This is experimental and not intended to replace the [Docker Compose](https://github.com/docker/compose) command-line tool. If you're looking to use Compose, head over to the [Compose installation instructions](http://docs.docker.com/compose/install/) to get started with it.**

```go
package main

import (
	"log"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
)

func main() {
	project, err := docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeFile: "docker-compose.yml",
			ProjectName: "my-compose",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	project.Up()
}
```

## Building

You need Docker and then run

    make


## Running

A partial implementation of the docker-compose CLI is also implemented in Go. The primary purpose of this code is so one can easily test the behavior of libcompose.

Run one of these:

```
docker-compose_darwin-386
docker-compose_linux-amd64
docker-compose_windows-amd64.exe
docker-compose_darwin-amd64
docker-compose_linux-arm
docker-compose_linux-386
docker-compose_windows-386.exe
```

### Tests

Make sure you have a full Go 1.4 environment with godeps

    godep go test ./...

This will be fully Dockerized in a bit

## Current status

The project is still being kickstarted... But it does a lot.  Please try it out and help us find bugs.

## Contributing

Want to hack on libcompose? [Docker's contributions guidelines](https://github.com/docker/docker/blob/master/CONTRIBUTING.md) apply.
