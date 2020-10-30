#!/bin/sh

first_install() {
  tmpdir=$(mktemp -d)

  echo "downloading files into ${tmpdir}"

  cd "${tmpdir}" && git clone https://github.com/loghole/tron.git 2>/dev/null

  latest=$(cd "${tmpdir}/tron" && git describe --tags --abbrev=0)

  cd "${tmpdir}/tron" && git checkout "${latest}" 2>/dev/null

  echo "build tron ${latest}"

  ldflags="-X 'github.com/loghole/tron/cmd/tron/internal/version.CliVersion=${latest}'"

  cd "${tmpdir}/tron/cmd/tron" && go build -o "$(go env GOPATH)/bin/tron" -ldflags "${ldflags}" main.go

  rm -rf "${tmpdir}"

  echo "tron ${latest} installed"
}

first_install
