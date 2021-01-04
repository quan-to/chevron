![Logo](assets/logo/chevron.png)

[![MIT License](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://tldrlegal.com/license/mit-license) [![Coverage Status](https://coveralls.io/repos/github/quan-to/chevron/badge.svg?branch=master)](https://coveralls.io/github/quan-to/chevron?branch=master) [![Build Status](https://travis-ci.org/quan-to/chevron.svg?branch=master)](https://travis-ci.org/quan-to/chevron)

A simple Toolkit to act as a GPG Creator / Signer / Verifier. This abstracts the use of the GPG and makes easy to sign / verify any GPG document using just a POST request.

# Usage

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

# Advanced Usage

* [Cluster Mode](https://github.com/quan-to/chevron/wiki/Cluster-Mode)
* [Vault Backend](https://github.com/quan-to/chevron/wiki/Hashicorp-Vault-Key-Backend)
* [Quanto Agent](https://github.com/quan-to/chevron/wiki/Quanto-Agent)
* [Binary Builds](https://github.com/quan-to/chevron/wiki/Binary-Builds)
* [Docker](https://github.com/quan-to/chevron/wiki/Docker)
* [Building](https://github.com/quan-to/chevron/wiki/Building)


Where is AgentUI??
==================

Agent-UI project has been moved to a separated repository. Check [https://github.com/quan-to/agent-ui](https://github.com/quan-to/agent-ui)

# Environment Variables

These are the Environment Variables that you can set to manage the webserver:

## Common Configuration

*   `PRIVATE_KEY_FOLDER` => Folder to load / store encrypted private keys. _(defaults to './keys')_
*   `MAX_KEYRING_CACHE_SIZE` => Maximum Number of Public Keys to cache (does not include Private Keys derived Public Keys). _(defaults to 1000)_
*   `SHOW_LINES` => Show filename and lines in logs
*   `REQUESTID_HEADER` => Header field to get request ID
*   `LOG_FORMAT` => Change log format (default is pipe delimited, provide the value `json` to log in JSON format)
*   `SKS_SERVER` => SKS Server to fetch / put public keys. _(defaults to 'http://pgp.mit.edu/')_
*   `KEY_PREFIX` => Prefix of the name of the keys to load (for example a key prefix `test_` will load any key named `test_XXXX`).
*   `MODE` => Mode of remote-signer (`single_key`, `default`)
*   `ON_DEMAND_KEY_LOAD` => Do not attempt to load all keys from keybackend. Load them as needed (defaults `false`)

## Caching Configuration

Remote Signer can use REDIS as a caching layer for GPG Keys and Tokens. If enabled, it also does some in-memory local caching with a smaller TTL.
To enable, use the following environment variables:

*   `REDIS_ENABLE` => `true` if should be enabled (`default: false`)
*   `REDIS_TLS_ENABLED` => `true` if TLS is enabled (`default: false`)
*   `REDIS_HOST` => Hostname of the REDIS server (`default: localhost:6379`)
*   `REDIS_USER` => Username of the REDIS server
*   `REDIS_PASS` => Password of the REDIS server
*   `REDIS_MAX_LOCAL_TTL` => Max local object TTL (in golang duration format): `default: 5m`
*   `REDIS_MAX_LOCAL_OBJECTS` => Max local objects (`default: 100`)
*   `REDIS_CLUSTER_MODE` => If the redis host is running in cluster mode. (`default: false`)

## Agent Configuration

*   `AGENT_TARGET_URL` => Target URL for Quanto Agent (defaults to `https://quanto-api.com.br/all`)
*   `AGENT_KEY_FINGERPRINT` => Default Key FingerPrint for Agent
*   `AGENT_BYPASS_LOGIN` => If the Login for using Quanto Agent should be bypassed. *DO NOT USE THIS IN EXPOSED REMOTESIGNER*
*   `AGENT_EXTERNAL_URL` => External URL used by GraphiQL to access agent. Defaults to `/agent`
*   `AGENTADMIN_EXTERNAL_URL` => External URL used by GraphiQL to access agent admin. Defaults to `/agentAdmin`
*   `READONLY_KEYPATH` => If the keypath is readonly. If `true` then it will create a temporary folder in `/tmp` and copy all keys to there so it can work over it. 
*   `HTTP_PORT` => HTTP Port that Remote Signer will run
*   Single Key Mode (`MODE=single_key`)
    * `SINGLE_KEY_PATH` => Path for the key to load as private key
    * `SINGLE_KEY_PASSWORD` => Password of the key to load as private key

## Cluster Mode Variables

*   `MASTER_GPG_KEY_PATH` => Master GPG Key Path
*   `MASTER_GPG_KEY_PASSWORD_PATH` => Master GPG Key Password Path
*   `MASTER_GPG_KEY_BASE64_ENCODED` => If the Master GPG Key is base64 encoded (default: true)
*   `SYSLOG_IP` => IP of the Syslog Server to send Console Messages _(defaults to '127.0.0.1')_ *Does not apply for Windows*
*   `SYSLOG_FACILITY` => Facility of the Syslog to use. _(defaults to 'LOG_USER')_
*   `DATABASE_DIALECT` => Dialect of the Database connection (`postgres`, `rethinkdb`. Defaults: none)
*   `CONNECTION_STRING` => Connection string for the database.

## Hashicorp Vault Key Backend Environment

*   `VAULT_STORAGE` => If a Hashicorp Vault should be used to store private keys instead of the disk (defaults `false`)
*   `VAULT_ADDRESS` => Hashicorp Vault URL
*   `VAULT_SKIP_VERIFY` => Hashicorp Vault Skip Verify SSL Certs on Connection
*   `VAULT_ROOT_TOKEN` => Hashicorp Vault Root Token
*   `VAULT_TOKEN_TTL` => Hashicorp Vault Token TTL (for example `24h`, default is `768h`. For more information see https://golang.org/pkg/time/#ParseDuration)
*   `VAULT_BACKEND` => Hashicorp Vault Backend (for example `secret`)
*   `VAULT_NAMESPACE` => if a Hashicorp Vault Namespace to use (appended to backend, for example if namespace is `remote-signer` the keys are stored under `secret/remote-signer`)

## Deprecated Environment Variables

**RethinkDB Usage is deprecated and discouraged**

*   `ENABLE_RETHINKDB_SKS` => Enables Internal SKS Server using RethinkDB (default: false)
    * Use `DATABASE_DIALECT=rethinkdb` instead
*   `RETHINK_TOKEN_MANAGER` => If a TokenManager using RethinkDB Should be used (defaults to `false`, uses MemoryTokenManager) [Requires ENABLE_RETHINK_SKS]
    * Use `DATABASE_TOKEN_MANAGER` instead
*   `RETHINK_AUTH_MANAGER` => If a AuthManager using RethinkDB Should be used (defaults to `false`, uses JSONAuthManager) [Requires ENABLE_RETHINK_SKS]
    * Use `DATABASE_AUTH_MANAGER` instead
*   `RETHINKDB_HOST` => Hostname of RethinkDB Server (default: "rethinkdb")
*   `RETHINKDB_USERNAME` => Username of RethinkDB Server (default "admin")
*   `RETHINKDB_PASSWORD` => Password of RethinKDB Server
*   `RETHINKDB_PORT` => Port of RethinkDB Server (default 28015)
*   `DATABASE_NAME` => RethinkDB Database Name (default "remote_signer")
