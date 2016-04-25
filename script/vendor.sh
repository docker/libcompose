#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."
rm -rf vendor/
source 'script/.vendor-helpers.sh'

clone git github.com/Sirupsen/logrus v0.9.0
clone git github.com/spf13/cobra 4c05eb1145f16d0e6bb4a3e1b6d769f4713cb41f
clone git github.com/spf13/pflag 8f6a28b0916586e7f22fe931ae2fcfc380b1c0e6
clone git github.com/docker/distribution d06d6d3b093302c02a93153ac7b06ebc0ffd1793
clone git github.com/vbatts/tar-split v0.9.11
clone git github.com/docker/docker v1.11.0
clone git github.com/docker/go-units 651fc226e7441360384da338d0fd37f2440ffbe3
clone git github.com/docker/go-connections v0.2.0
clone git github.com/docker/engine-api a6dca654f28f26b648115649f6382252ada81119
clone git github.com/flynn/go-shlex 3f9db97f856818214da2e1057f8ad84803971cff
clone git github.com/gorilla/context 14f550f51a
clone git github.com/gorilla/mux e444e69cbd
clone git github.com/opencontainers/runc 7b6c4c418d5090f4f11eee949fdf49afd15838c9
clone git github.com/stretchr/testify a1f97990ddc16022ec7610326dd9bce31332c116
clone git github.com/davecgh/go-spew 5215b55f46b2b919f50a1df0eaa5886afe4e3b3d
clone git github.com/pmezard/go-difflib d8ed2627bdf02c080bf22230dbb337003b7aba2d
clone git golang.org/x/crypto 4d48e5fa3d62b5e6e71260571bf76c767198ca02 https://github.com/golang/crypto.git
clone git golang.org/x/net 47990a1ba55743e6ef1affd3a14e5bac8553615d https://github.com/golang/net.git
clone git gopkg.in/check.v1 11d3bc7aa68e238947792f30573146a3231fc0f1
clone git github.com/Azure/go-ansiterm 70b2c90b260171e829f1ebd7c17f600c11858dbe
clone git github.com/cloudfoundry-incubator/candiedyaml 5cef21e2e4f0fd147973b558d4db7395176bcd95
clone git github.com/Microsoft/go-winio v0.1.0
clone git github.com/xeipuuv/gojsonpointer e0fe6f68307607d540ed8eac07a342c33fa1b54a
clone git github.com/xeipuuv/gojsonreference e02fc20de94c78484cd5ffb007f8af96be030a45
clone git github.com/xeipuuv/gojsonschema ac452913faa25c08bb78810d3e6f88b8a39f8f25
clone git github.com/kr/pty 5cf931ef8f

clean && mv vendor/src/* vendor
