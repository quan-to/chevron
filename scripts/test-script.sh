#!/bin/bash

# Assumes racerxdl/goxenialtest and repository mounted to /go/src/github.com/quan-to/chevron
set -e

export DEBIAN_FRONTEND=noninteractive
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_DEV_ROOT_TOKEN_ID="123456"
export VAULT_ROOT_TOKEN=$VAULT_DEV_ROOT_TOKEN_ID
export VAULT_ADDRESS="$VAULT_ADDR"
export VAULT_USE_USERPASS="true"
export VAULT_USERNAME="remotesigner"
export VAULT_PASSWORD="123456"
export DATABASE_NAME="qrs_test"

echo "----- Installing dependencies -------"


echo "deb http://download.rethinkdb.com/apt xenial main" | tee /etc/apt/sources.list.d/rethinkdb.list
wget -qO- https://download.rethinkdb.com/apt/pubkey.gpg | apt-key add -

apt-get -qqy update
apt install -y --no-install-recommends rethinkdb


go get github.com/mattn/goveralls
wget https://releases.hashicorp.com/vault/1.0.2/vault_1.0.2_linux_amd64.zip
unzip vault_1.0.2_linux_amd64.zip
cp vault /usr/bin

echo "------ Starting Services -------"
echo "Starting RethinkDB"
rethinkdb 2>> /dev/null 1>>/dev/null & echo $! > $HOME/rethinkdb.pid
echo "RethinkDB PID `cat $HOME/rethinkdb.pid`"

echo "Starting Vault"
vault server -dev 2>> /dev/null 1>>/dev/null & echo $! > $HOME/vault.pid
echo "Waiting vault settle" & sleep 2
echo "Vault started with PID `cat $HOME/vault.pid`"

cd /go/src/github.com/quan-to/chevron/

echo "Adding test-policy.hcl to vault"
vault auth enable userpass
vault policy write test-policy ./test-policy.hcl
vault write auth/userpass/users/${VAULT_USERNAME} password=${VAULT_PASSWORD} policies=test-policy

echo "Downloading golangci-lint"
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $GOPATH/bin v1.10.2

echo "----- Running go get ------"
go get
for i in cmd/*
do
  cd $i
  go get
  cd ../..
done

echo "------ Running Tests ------"

golangci-lint run
go test -v -race ./... -coverprofile=qrs.coverprofile
goveralls -coverprofile=qrs.coverprofile

echo "----- Closing Services -----"
echo "Closing RethinkDB"
kill -9 `cat ~/rethinkdb.pid`
echo "Closing Vault"
kill -9 `cat ~/vault.pid`
echo "Done"
