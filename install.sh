#!/bin/sh

first_install() {
  tmpdir=$(mktemp -d)

  echo "downloading files into ${tmpdir}"

  cd "${tmpdir}" && git clone https://github.com/loghole/tron.git

  latest=$(cd "${tmpdir}" && git describe --tags --always | perl -pe 's/cmd\/tron\///')

  cd "${tmpdir}" && git checkout "$(latest)"

  cd "${tmpdir}/tron/cmd/tron" && make build

  rm -rf "${tmpdir}"
}

first_install
