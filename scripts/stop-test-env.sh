#!/bin/bash

echo "Closing Vault"
kill -9 `cat ~/vault.pid`
rm ~/vault.pid