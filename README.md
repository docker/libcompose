# libcompose

Go library for compose and full docker-compose CLI implementation in go.

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

The project is still being kickstarted... But it does a lot.

