#!/usr/bin/env bash
#
# On a Debian/Ubuntu system, bootstraps a docker install and the cockroach
# repo.

set -euxo pipefail

curl -fsSL https://deb.nodesource.com/gpgkey/nodesource.gpg.key | sudo apt-key add -
echo "deb https://deb.nodesource.com/node_6.x xenial main" | sudo tee /etc/apt/sources.list.d/nodesource.list

curl -fsSL https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list

sudo apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get dist-upgrade -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold"
sudo apt-get install -y --no-install-recommends \
  mosh \
  autoconf \
  cmake \
  ccache \
  docker.io \
  libncurses-dev \
  make \
  gcc \
  g++ \
  git \
  nodejs \
  yarn \
  bison

sudo adduser "${USER}" docker

# Configure environment variables.
echo 'export PATH="/usr/lib/ccache:${PATH}:/usr/local/go/bin"' >> ~/.bashrc_bootstrap
echo 'export COCKROACH_BUILDER_CCACHE=1' >> ~/.bashrc_bootstrap
echo '. ~/.bashrc_bootstrap' >> ~/.bashrc
. ~/.bashrc_bootstrap

# Install Go.
trap 'rm -f /tmp/go.tgz' EXIT
curl https://dl.google.com/go/go1.11.6.linux-amd64.tar.gz > /tmp/go.tgz
sha256sum -c - <<EOF
4e1864282d8d20010d6385a12a1e35641783a380a7c57907bfb46a5499c5eb49  /tmp/go.tgz
EOF
sudo tar -C /usr/local -zxf /tmp/go.tgz

# Clone CockroachDB.
git clone https://github.com/weisslj/cockroach "$(go env GOPATH)/src/github.com/weisslj/cockroach"

# Install the Unison file-syncer.
. bootstrap/bootstrap-unison.sh
