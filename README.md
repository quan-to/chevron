![Logo](https://raw.githubusercontent.com/quan-to/chevron/develop/logo/chevron.png)

[![MIT License](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://tldrlegal.com/license/mit-license) [![Coverage Status](https://coveralls.io/repos/github/quan-to/chevron/badge.svg?branch=master)](https://coveralls.io/github/quan-to/chevron?branch=master) [![Build Status](https://travis-ci.org/quan-to/chevron.svg?branch=master)](https://travis-ci.org/quan-to/chevron)

A simple Toolkit to act as a GPG Creator / Signer / Verifier. This abstracts the use of the GPG and makes easy to sign / verify any GPG document using just a POST request.

Usage
=======

* [Getting Started](https://github.com/quan-to/chevron/wiki)
* [Creating GPG Keys](https://github.com/quan-to/chevron/wiki/Creating-GPG-keys)
* [Setting up Keys](https://github.com/quan-to/remote-signer/wiki/SettingUp-keys)
* [Listing loaded private keys](https://github.com/quan-to/chevron/wiki/List-loaded-private-keys)
* [Unlock Private Key](https://github.com/quan-to/remote-signer/wiki/Unlock-private-key)
* [Signing Data](https://github.com/quan-to/remote-signer/wiki/Signing-Data)
* [Listing cached public keys](https://github.com/quan-to/chevron/wiki/List-cached-public-keys)
* [Verifying Signatures](https://github.com/quan-to/remote-signer/wiki/Verifying-Signatures)
* [Encrypting Data](https://github.com/quan-to/remote-signer/wiki/Encrypting-Data)
* [Decrypting Data](https://github.com/quan-to/remote-signer/wiki/Decrypting-Data)

Advanced Usage
==============

* [Cluster Mode](https://github.com/quan-to/chevron/wiki/Cluster-Mode)
* [Vault Backend](https://github.com/quan-to/chevron/wiki/Hashicorp-Vault-Key-Backend)
* [Quanto Agent](https://github.com/quan-to/chevron/wiki/Quanto-Agent)
* Binary Builds
* Docker
* [Building](https://github.com/quan-to/chevron/wiki/Building)


Environment Variables
=====================

These are the Environment Variables that you can set to manage the webserver:

*   `PRIVATE_KEY_FOLDER` => Folder to load / store encrypted private keys. _(defaults to './keys')_
*   `SYSLOG_IP` => IP of the Syslog Server to send Console Messages _(defaults to '127.0.0.1')_ *Does not apply for Windows*
*   `SYSLOG_FACILITY` => Facility of the Syslog to use. _(defaults to 'LOG_USER')_
*   `SKS_SERVER` => SKS Server to fetch / put public keys. _(defaults to 'http://pgp.mit.edu/')_
*   `KEY_PREFIX` => Prefix of the name of the keys to load (for example a key prefix `test_` will load any key named `test_XXXX`).
*   `MAX_KEYRING_CACHE_SIZE` => Maximum Number of Public Keys to cache (does not include Private Keys derived Public Keys). _(defaults to 1000)_
*   `ENABLE_RETHINKDB_SKS` => Enables Internal SKS Server using RethinkDB (default: false)
*   `RETHINKDB_HOST` => Hostname of RethinkDB Server (default: "rethinkdb")
*   `RETHINKDB_USERNAME` => Username of RethinkDB Server (default "admin")
*   `RETHINKDB_PASSWORD` => Password of RethinKDB Server
*   `RETHINK_TOKEN_MANAGER` => If a TokenManager using RethinkDB Should be used (defaults to `false`, uses MemoryTokenManager) [Requires ENABLE_RETHINK_SKS]
*   `RETHINK_AUTH_MANAGER` => If a AuthManager using RethinkDB Should be used (defaults to `false`, uses JSONAuthManager) [Requires ENABLE_RETHINK_SKS]
*   `RETHINKDB_PORT` => Port of RethinkDB Server (default 28015)
*   `AGENT_TARGET_URL` => Target URL for Quanto Agent (defaults to `https://quanto-api.com.br/all`)
*   `AGENT_KEY_FINGERPRINT` => Default Key FingerPrint for Agent
*   `AGENT_BYPASS_LOGIN` => If the Login for using Quanto Agent should be bypassed. *DO NOT USE THIS IN EXPOSED REMOTESIGNER*
*   `AGENT_EXTERNAL_URL` => External URL used by GraphiQL to access agent. Defaults to `/agent`
*   `AGENTADMIN_EXTERNAL_URL` => External URL used by GraphiQL to access agent admin. Defaults to `/agentAdmin`
*   `DATABASE_NAME` => RethinkDB Database Name (default "remote_signer")
*   `MASTER_GPG_KEY_PATH` => Master GPG Key Path
*   `MASTER_GPG_KEY_PASSWORD_PATH` => Master GPG Key Password Path
*   `MASTER_GPG_KEY_BASE64_ENCODED` => If the Master GPG Key is base64 encoded (default: true)
*   `VAULT_ADDRESS` => Hashicorp Vault URL
*   `VAULT_SKIP_VERIFY` => Hashicorp Vault Skip Verify SSL Certs on Connection
*   `VAULT_ROOT_TOKEN` => Hashicorp Vault Root Token
*   `VAULT_BACKEND` => Hashicorp Vault Backend (for example `secret`)
*   `VAULT_STORAGE` => If a Hashicorp Vault should be used to store private keys instead of the disk
*   `VAULT_NAMESPACE` => if a Hashicorp Vault Namespace to use (appended to backend, for example if namespace is `remote-signer` the keys are stored under `secret/remote-signer`)
*   `HTTP_PORT` => HTTP Port that Remote Signer will run
*   `READONLY_KEYPATH` => If the keypath is readonly. If `true` then it will create a temporary folder in `/tmp` and copy all keys to there so it can work over it. 

