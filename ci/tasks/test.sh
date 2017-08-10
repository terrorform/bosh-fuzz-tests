#!/usr/bin/env bash

set -e

cp $(ls $CLI_DIR_PATH/alpha-bosh-cli-*-linux-amd64) "/tmp/gobosh"
chmod a+x "/tmp/gobosh"

export PATH=/usr/local/ruby/bin:/usr/local/go/bin:$PATH
export DB='postgresql'

echo 'Starting DB...'
su postgres -c '
  export PATH=/usr/lib/postgresql/9.4/bin:$PATH
  export PGDATA=/tmp/postgres
  export PGLOGS=/tmp/log/postgres
  mkdir -p $PGDATA
  mkdir -p $PGLOGS
  initdb -U postgres -D $PGDATA
  pg_ctl start -l $PGLOGS/server.log
'

source /etc/profile.d/chruby.sh
chruby $RUBY_VERSION

bosh_src_path="$PWD/$BOSH_SRC_PATH"

echo 'Installing dependencies...'

gem install cf-uaac --no-document

agent_path=bosh-src/src/go/src/github.com/cloudfoundry/
mkdir -p $agent_path
cp -r bosh-agent $agent_path

(
  cd $bosh_src_path
  bundle install --local
  bundle exec rake spec:integration:install_dependencies

  echo "Building agent..."
  go/src/github.com/cloudfoundry/bosh-agent/bin/build
)

echo 'Running tests...'

export GOPATH=$(realpath bosh-fuzz-tests)

sed -i s#BOSH_SRC_PATH#${bosh_src_path}#g bosh-fuzz-tests/ci/concourse-config.json
sed -i s#PWD#${PWD}#g bosh-fuzz-tests/ci/concourse-config.json

cp bosh-fuzz-tests/assets/ssl/* /tmp/

pushd bosh-fuzz-tests/src/github.com/cloudfoundry-incubator/bosh-fuzz-tests/
	go install github.com/onsi/ginkgo/ginkgo
	ginkgo -r -p -randomizeAllSpecs -randomizeSuites -race
popd

go run bosh-fuzz-tests/src/github.com/cloudfoundry-incubator/bosh-fuzz-tests/main.go bosh-fuzz-tests/ci/concourse-config.json ${SEED:-}
