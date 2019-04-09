#!/bin/bash

echo "Starting Vault"
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_DEV_ROOT_TOKEN_ID="123456"
export VAULT_ROOT_TOKEN=$VAULT_DEV_ROOT_TOKEN_ID
export VAULT_ADDRESS="$VAULT_ADDR"
export VAULT_USE_USERPASS="true"
export VAULT_USERNAME="remotesigner"
export VAULT_PASSWORD="123456"
export DATABASE_NAME="qrs_test"
export DO_START_RETHINK="true"
vault server -dev 2>&1 1>vault.log & echo $! > $HOME/vault.pid
echo "Waiting vault settle" & sleep 2
echo "Vault started with PID `cat $HOME/vault.pid`"
vault auth enable userpass
vault policy write test-policy ./test-policy.hcl
vault write auth/userpass/users/${VAULT_USERNAME} password=${VAULT_PASSWORD} policies=test-policy

echo "Done. Please run this before the tests:
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_DEV_ROOT_TOKEN_ID="123456"
export VAULT_ROOT_TOKEN=$VAULT_DEV_ROOT_TOKEN_ID
export VAULT_ADDRESS=\"$VAULT_ADDR\"
export VAULT_USE_USERPASS=\"true\"
export VAULT_USERNAME=\"remotesigner\"
export VAULT_PASSWORD=\"123456\"
export DATABASE_NAME=\"qrs_test\"
export DO_START_RETHINK=\"true\"
"