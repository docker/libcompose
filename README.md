### :warning: Deprecation Notice: This project and repository is now deprecated and is no longer under active development. Please use [compose-go](https://github.com/compose-spec/compose-go) instead.

# libcompose

[![GoDoc](https://godoc.org/github.com/docker/libcompose?status.png)](https://godoc.org/github.com/docker/libcompose)
[![Build Status](https://jenkins.dockerproject.org/job/docker/job/libcompose/branch/master/badge/icon)](https://jenkins.dockerproject.org/job/docker/job/libcompose/branch/master/)

A Go library for Docker Compose. It does everything the command-line tool does, but from within Go -- read Compose files, start them, scale them, etc.

**Note: This is not really maintained anymore — the reason are diverse but mainly lack of time from the maintainers**

The current state is the following :
- The `libcompose` CLI should considered abandonned. The `v2` parsing is incomplete and `v3` parsing is missing.
- The official compose Go parser implementation is on [`docker/cli`](https://github.com/docker/cli/tree/master/cli/compose) but only support `v3` version of the compose format.

What is the work that is needed:
- Remove the cli code (thus removing dependencies to `docker/cli` )
- Clearer separation of packages : `parsing`, `conversion` (to docker api or swarm api), `execution` (`Up`, `Down`, … behaviors)
- Add support for all compose format version (v1, v2.x, v3.x)
- Switch to either `golang/dep` or `go mod` for dependencies (removing the `vendor` folder)
- *(bonus)* extract the [`docker/cli`](https://github.com/docker/cli/tree/master/cli/compose) code here and vendor this library into `docker/cli`.

If you are interested to work on `libcompose`, feel free to ping me (over twitter @vdemeest), I'll definitely do code reviews and help as much as I can 😉.

**Note: This is experimental and not intended to replace the [Docker Compose](https://github.com/docker/compose) command-line tool. If you're looking to use Compose, head over to the [Compose installation instructions](http://docs.docker.com/compose/install/) to get started with it.**

Here is a list of known project that uses `libcompose`:

- [rancher-compose](https://github.com/rancher/rancher-compose) and [rancher os](https://github.com/rancher/os) (by [Rancher](https://github.com/rancher))
- [openshift](https://github.com/openshift/origin) (by [Red Hat](https://github.com/openshift))
- [henge](https://github.com/redhat-developer/henge) (by [Red Hat](https://github.com/redhat-developer)) [Deprecated in favour of kompose]
- [kompose](https://github.com/skippbox/kompose) (by [skippbox](https://github.com/skippbox))
- [compose2kube](https://github.com/kelseyhightower/compose2kube) (by [kelseyhightower](https://github.com/kelseyhightower))
- [amazon-ecs-cli](https://github.com/aws/amazon-ecs-cli) (by [Amazon AWS](https://github.com/aws))
- [libkermit](https://github.com/libkermit/docker) (by [vdemeester](https://github.com/vdemeester))

## Usage

```go
package main

import (
	"log"

	"golang.org/x/net/context"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

func main() {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			ProjectName:  "my-compose",
		},
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Up(context.Background(), options.Up{})

	if err != nil {
		log.Fatal(err)
	}
}
```

## Tests (unit & integration)


You can run unit tests using the `test-unit` target and the
integration test using the `test-integration` target. If you don't use
Docker and `make` to build `libcompose`, you can use `go test` and the
following scripts : `hack/test-unit` and `hack/test-integration`.

```bash
$ make test-unit
docker build -t "libcompose-dev:refactor-makefile" .
#[…]
---> Making bundle: test-unit (in .)
+ go test -cover -coverprofile=cover.out ./docker
ok      github.com/docker/libcompose/docker     0.019s  coverage: 4.6% of statements
+ go test -cover -coverprofile=cover.out ./project
ok      github.com/docker/libcompose/project    0.010s  coverage: 8.4% of statements
+ go test -cover -coverprofile=cover.out ./version
ok      github.com/docker/libcompose/version    0.002s  coverage: 0.0% of statements

Test success
```


## Current status

The project is still being kickstarted... But it does a lot.  Please try it out and help us find bugs.

## Contributing

Want to hack on libcompose? [Docker's contributions guidelines](https://github.com/docker/libcompose/blob/master/CONTRIBUTING.md) apply.

If you have comments, questions, or want to use your knowledge to help other, come join the conversation on IRC. You can reach us at #libcompose on Freenode.
