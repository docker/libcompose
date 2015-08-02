# libcompose

An official implementation of Docker Compose in Go made available as a library.

**Note: This is an experimental alternate implementation of [Docker Compose](https://github.com/docker/compose)**

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

    ./script/build


## Running

A full implementation of the docker-compose CLI is implemented also in Go.  The primary purpose of this code is to provide a way in which one can easily test the correctness of the behavior of libcompose.

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

