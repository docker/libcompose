#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."
rm -rf vendor/
source 'script/.vendor-helpers.sh'

clone git github.com/Sirupsen/logrus v0.8.2
clone git github.com/codegangsta/cli 6086d7927ec35315964d9fea46df6c04e6d697c1
clone git github.com/docker/distribution c6c9194e9c6097f84b0ff468a741086ff7704aa3
clone git github.com/docker/docker 58b270c338e831ac6668a29788c72d202f9fc251
clone git github.com/docker/libtrust 9cbd2a1374f46905c68a4eb3694a130610adc62a
clone git github.com/flynn/go-shlex 3f9db97f856818214da2e1057f8ad84803971cff
clone git github.com/fsouza/go-dockerclient 39d9fefa6a7fd4ef5a4a02c5f566cb83b73c7293
clone git github.com/gorilla/context 215affda49addc4c8ef7e2534915df2c8c35c6cd
clone git github.com/gorilla/mux f15e0c49460fd49eebe2bcc8486b05d1bef68d3a
clone git github.com/opencontainers/runc b40c7901845dcec5950ecb37cb9de178fc2c0604
clone git github.com/stretchr/testify 7e4a149930b09fe4c2b134c50ce637457ba6e966
clone git golang.org/x/crypto 4d48e5fa3d62b5e6e71260571bf76c767198ca02 https://github.com/golang/crypto.git
clone git golang.org/x/net 3a29182c25eeabbaaf94daaeecbc7823d86261e7 https://github.com/golang/net.git
clone git gopkg.in/check.v1 11d3bc7aa68e238947792f30573146a3231fc0f1
clone git gopkg.in/yaml.v2 49c95bdc21843256fb6c4e0d370a05f24a0bf213
clone git github.com/Azure/go-ansiterm 70b2c90b260171e829f1ebd7c17f600c11858dbe
clone git github.com/cloudfoundry-incubator/candiedyaml 55a459c2d9da2b078f0725e5fb324823b2c71702

clean && mv vendor/src/* vendor
