package app

import (
	"github.com/codegangsta/cli"
	"github.com/docker/libcompose/docker"
)

func DockerClientFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "tls",
			Usage: "Use TLS; implied by --tlsverify",
		},
		cli.BoolFlag{
			Name:  "tlsverify",
			Usage: "Use TLS and verify the remote",
		},
		cli.StringFlag{
			Name:  "tlscacert",
			Usage: "Trust certs signed only by this CA",
		},
		cli.StringFlag{
			Name:  "tlscert",
			Usage: "Path to TLS certificate file",
		},
		cli.StringFlag{
			Name:  "tlskey",
			Usage: "Path to TLS key file",
		},
		cli.StringFlag{
			Name:  "configdir",
			Usage: "Path to docker config dir, default ${HOME}/.docker",
		},
	}
}

func Populate(context *docker.Context, c *cli.Context) {
	context.Tls = c.Bool("tls")
	context.TlsVerify = c.Bool("tlsverify")
	context.Ca = c.String("tlscacert")
	context.Cert = c.String("tlscert")
	context.Key = c.String("tlskey")
	context.ConfigDir = c.String("configdir")
}
