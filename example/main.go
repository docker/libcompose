package main

import (
	"log"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

func main() {
	project, err := docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			ProjectName:  "yeah-compose",
		}
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Up(options.Up{})

	if err != nil {
		log.Fatal(err)
	}
}
