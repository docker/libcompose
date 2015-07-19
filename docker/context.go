package docker

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/cliconfig"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/homedir"
	"github.com/docker/libcompose/project"
	"github.com/samalba/dockerclient"
)

const (
	defaultTrustKeyFile = "key.json"
	defaultCaFile       = "ca.pem"
	defaultKeyFile      = "key.pem"
	defaultCertFile     = "cert.pem"
)

var (
	dockerCertPath  = os.Getenv("DOCKER_CERT_PATH")
	dockerTlsVerify = os.Getenv("DOCKER_TLS_VERIFY") != ""
)

type Context struct {
	project.Context
	Builder    Builder
	Client     dockerclient.Client
	Tls        bool
	TlsVerify  bool
	TrustKey   string
	Ca         string
	Cert       string
	Key        string
	Host       string
	ConfigDir  string
	ConfigFile *cliconfig.ConfigFile
	tlsConfig  *tls.Config
}

func (c *Context) open() error {
	err := c.LookupConfig()
	if err != nil {
		return err
	}

	return c.CreateClient()
}

func (c *Context) LookupConfig() error {
	if c.ConfigFile != nil {
		return nil
	}

	config, err := cliconfig.Load(c.ConfigDir)
	if err != nil {
		return err
	}

	c.ConfigFile = config

	return nil
}

func (c *Context) CreateClient() error {
	if c.Client != nil {
		return nil
	}

	if c.Ca == "" {
		c.Ca = filepath.Join(dockerCertPath, defaultCaFile)
	}
	if c.Cert == "" {
		c.Cert = filepath.Join(dockerCertPath, defaultCertFile)
	}
	if c.Key == "" {
		c.Key = filepath.Join(dockerCertPath, defaultKeyFile)
	}

	if c.Host == "" {
		defaultHost := os.Getenv("DOCKER_HOST")
		if defaultHost == "" {
			if runtime.GOOS != "windows" {
				// If we do not have a host, default to unix socket
				defaultHost = fmt.Sprintf("unix://%s", opts.DefaultUnixSocket)
			} else {
				// If we do not have a host, default to TCP socket on Windows
				defaultHost = fmt.Sprintf("tcp://%s:%d", opts.DefaultHTTPHost, opts.DefaultHTTPPort)
			}
		}
		defaultHost, err := opts.ValidateHost(defaultHost)
		if err != nil {
			return err
		}
		c.Host = defaultHost
	}

	if c.TrustKey == "" {
		c.TrustKey = filepath.Join(homedir.Get(), ".docker", defaultTrustKeyFile)
	}

	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true

	// Regardless of whether the user sets it to true or false, if they
	// specify --tlsverify at all then we need to turn on tls
	if c.TlsVerify {
		c.Tls = true
	}

	// If we should verify the server, we need to load a trusted ca
	if c.TlsVerify {
		certPool := x509.NewCertPool()
		file, err := ioutil.ReadFile(c.Ca)
		if err != nil {
			logrus.Errorf("Couldn't read ca cert %s: %s", c.Ca, err)
			return err
		}
		certPool.AppendCertsFromPEM(file)
		tlsConfig.RootCAs = certPool
		tlsConfig.InsecureSkipVerify = false
	}

	// If tls is enabled, try to load and send client certificates
	if c.Tls {
		_, errCert := os.Stat(c.Cert)
		_, errKey := os.Stat(c.Key)
		if errCert == nil && errKey == nil {
			c.Tls = true
			cert, err := tls.LoadX509KeyPair(c.Cert, c.Key)
			if err != nil {
				logrus.Errorf("Couldn't load X509 key pair: %q. Make sure the key is encrypted", err)
				return err
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
		// Avoid fallback to SSL protocols < TLS1.0
		tlsConfig.MinVersion = tls.VersionTLS10
	}

	if c.Tls {
		c.tlsConfig = &tlsConfig
	}

	client, err := dockerclient.NewDockerClient(c.Host, c.tlsConfig)
	c.Client = client
	return err
}
