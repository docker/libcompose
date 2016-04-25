package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/spf13/cobra"
)

// Populate updates the specified docker context based on command line arguments and subcommands.
func Populate(context *docker.Context, c *cobra.Command) {
	context.ConfigDir, _ = c.Flags().GetString("configdir")

	opts := docker.ClientOpts{}
	opts.TLS, _ = c.Flags().GetBool("tls")
	opts.TLSVerify, _ = c.Flags().GetBool("tlsverify")
	opts.TLSOptions.CAFile, _ = c.Flags().GetString("tlscacert")
	opts.TLSOptions.CertFile, _ = c.Flags().GetString("tlscert")
	opts.TLSOptions.KeyFile, _ = c.Flags().GetString("tlskey")

	clientFactory, err := docker.NewDefaultClientFactory(opts)
	if err != nil {
		logrus.Fatalf("Failed to construct Docker client: %v", err)
	}

	context.ClientFactory = clientFactory
}
