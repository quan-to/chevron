Quanto Remote Signer (QRS)
====================

[![MIT License](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://tldrlegal.com/license/mit-license) [![Coverage Status](https://coveralls.io/repos/github/quan-to/remote-signer/badge.svg?branch=GoLang)](https://coveralls.io/github/quan-to/remote-signer?branch=master) [![Build Status](https://travis-ci.org/quan-to/remote-signer.svg?branch=master)](https://travis-ci.org/quan-to/remote-signer)

A simple Web Server to act as a GPG Creator / Signer / Verifier. This abstracts the use of the GPG and makes easy to sign / verify any GPG document using just a POST request.

Please notice that this application is *NOT inteded to ran public in the internet*. This is inteded to be a helper service to your application be able to sign / verify data (same as local gpg in the system). Because of that, it only listens for localhost.


TODO
====

*   Increment this document with models and enums

Usage
=====

This application opens up a WebServer listening on port *5100* and have a base URL defined as `/remoteSigner`.

#### Setting up GPG Private Keys

By default QRS searchs for encrypted private keys at `./keys`. Put all the private keys you want to use in Encrypted Ascii Armored Format inside it. It will iterate over all files and load them. If you don't have one, you can either create using the `gpg` toolkit or by calling the create api. Notice that calling the create API does not automatically store the key at the `keys` folder.

The keys folder can be overrided by the `PRIVATE_KEY_FOLDER` environment variable.

#### Creating a GPG Key

Although creating a GPG Key here might not be a good idea, you can use QRS to generate new GPG Keys on the fly. To do so, make a POST request to `/remoteSigner/gpg/generateKey` with the following JSON Content:

```json
{
  "Identifier": "Lucas Teske <lucas@teske.com.br>",
  "Password": "123456",
  "Bits": 3072
}
```

It should return your Encrypted GPG Private Key in ASCII Armored format.

```
-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: BCPG C# v1.8.1.0

lQVsBFpqVT0BDADIMVGd96DMUGf+zrcs0cGTzofbvV56WTWFju9WzIiUMigON6Qw
XdpdHUad1H31pnI1COCKH+k2t3TOlQr7qgXHMFOjW+/xHKoN6NhGMZVC7MkUllaj
uTFDH9823N/fhbJ4BRuBb2a5X4HBIeIDscu19xsW5B3HvwggojjhZ5iKRCt49Hsv
dJ6gPA5fDURGAbt9xdAqWvlkT9xagHqylVSG1A1CxOmeP3p+Vfjh/IhCgZ/nbi52
s+iBthuraYJAIPB9snASniMIqYs7sWTpC8T4m+WYEZGB2ejvVscmEgXFNWn6hzKI
(...)
-----END PGP PRIVATE KEY BLOCK-----
```

#### Unlocking a Private Key
Before any sign operation can be done, you need to decrypt the loaded private keys. If the keys is stored in a non-encrypted format (no password) you don't need to do that step. Simple call `/remoteSigner/gpg/unlockKey` with the following JSON payload:

```json
{
  "FingerPrint": "D7362B4CC546DB11",
  "Password": "123456"
}
```

#### Signing Data
Please check https://github.com/quan-to/remote-signer/wiki/Sign_data

#### Verifing Signatures
Please check https://github.com/quan-to/remote-signer/wiki/Verifying_Signatures

#### Adding Encrypted Private Key through API

To add a private key, you can make a POST to `/remoteSigner/keyRing/addPrivateKey` with the following payload:

```json
{
  "EncryptedPrivateKey": "-----BEGIN PGP PRIVATE KEY BLOCK-----\nVersion: GnuPG (...) JSlmyLSuTHXzeKo72hP40y3Xkf\nuugqVOWHeE7v7ARMu1mhXS6qWzZmxsjixV1d0kXSo9LzUyFqNtkasUiL2aoXQ70z\nlbMia0X7KJbYnbG5XLEDiMjzDQ==\n=JwWd\n-----END PGP PRIVATE KEY BLOCK-----",
  "SaveToDisk": true
}
```

The `SaveToDisk` parameter tells the server to save that private keys in the `KeysFolder`.

#### List Cached Public Keys

Execute GET to `/remoteSigner/keyRing/cachedKeys`

Returns:

```json
[
    {
        "FingerPrint": "D7362B4CC546DB11",
        "Identifier":"Benchmark Test Key",
        "Bits": 4096,
        "ContainsPrivateKey": false,
        "PrivateKeyDecrypted": false
    }
]
```

#### List Loaded Private Keys
Execute GET to `/remoteSigner/keyRing/privateKeys`

Returns:

```json
[
    {
        "FingerPrint": "D7362B4CC546DB11",
        "Identifier": "Benchmark Test Key",
        "Bits": 4096,
        "ContainsPrivateKey": true,
        "PrivateKeyDecrypted": false
    }
]
```

#### Encrypt Data using public key
Please check https://github.com/quan-to/remote-signer/wiki/Encrypt_data

#### Decrypt Data using decrypted private key
Ensure that the private key from the data you're trying to decrypt is loaded and decrypted, then execute a POST to `/remoteSigner/gpg/decrypt` with the following payload:

```json
{
  "AsciiArmoredData": "-----BEGIN PGP MESSAGE-----\nVersion: BCPG C# v1.8.1.0\n\nhQILAwAWqcqHCvpZAQ/4vF53gHVus8aKyKGkzb7jn2R4aZB3KCQ08S2xhAUvZFF8\n0qaeLxPGDdOo4X43zNmOvfIth4IwnDFF/SlD6E9ToxI+oDBC2hU92GyQZmlrb0dm\nHfVtKxCP9D6bAUHb9/G2QrbLSwov7TKlYs4gcqv72Lh4It8wVZaUm+qWb1EL4I33\nM+RHPSkmDPpCVWJxPQv/5Bt0h48wX9V7JtFc2FXJgJhYyrRxxIEFDcof4jdbvH/2\ncd5DDoLJPq3w4R4GKLxgQisgK2fp9jsl5AUzBiNy++l80rJW3m9TL3hLHUqqvL2R\nZrglp5KCR367uB0b6H+oXCRkxsgulTtXWM111HfTEJ0FYkYMfjxwYLdgjeduilhP\n7bkLDXFJgR3TaxHweUx4tOYRREsRSnzlDEDt+RdCHnP27mn/8QOi2wzi5zTP/KIr\nHlNfb+yw1BlFS5swFp+QLj7/QfkZefsneQC+zKfzyV9Hyz0b5tqXmsn+aVREF9D7\nQpqbEyHO1E/amz0hPIqu8CIIxr9Exmjxj5jV4MRgqVZ+5ukjiahG4jnnGuMPMvlp\nYjdZ1lAq0LDs+XSf9QbZE63j4YDT5tuJXNhUYojhb+DSSlL5LmQCuzTtZLNZaS+S\n26fQ5R3NYTWsJF3gvqwyXCr/49gYDxU2YNOBGdHGvOsDqnqchceRqWXfCJ8QlNKG\nAXLylUqQH0y58X0DTbEUEDtKRHAk42f9hicpxQY0FfnrUnggIBFubs385k6LIIDR\n4Xs6LwwjGFT9XqWzNa7adi+60sfrlN2iTYRZJGsNdvGmnMTClS0e6i6rlgQJAqHe\nl5rf7WGniCF+sAjxmbJ53TPBrh/sUlMMl0acXmXz4EZxnaENBJg=\n=n4Ks\n-----END PGP MESSAGE-----",
}
```

Returns:

```json
{
  "FingerPrint": "0016A9CA870AFA59",
  "Base64Data": "eyJxdWVyeSI6InF1ZXJ5IHsgR2V0QmFua1N5c3RlbVN0YXR1cyAoYmFua051bWJlcjogXCI2MzNcIikgfSJ9Cg==",
  "Filename": "QuantoEncrypt-1518071090398.bin",
  "IsIntegrityProtected": true,
  "IsIntegrityOK": true
}
```

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

Cluster Mode
============

TODO

Vault Backend
=============

TODO

Quanto Agent
============

TODO

Binary Builds
=============

TODO

Docker
======

TODO

Building
========

# Prepare

First of all, you should be aware how to configure your golang environment (which should be similar to all Operating Systems). For a more specific how-to please refer to official golang install: https://golang.org/doc/install

# Building (any os)

Since Remote Signer is a pure golang program, its build instructions are the same for *any* operating system. 
```bash
cd cmd/server
go build -o remote-signer
```

If you're on windows, run instead

```powershell
cd cmd/server
go build -o remote-signer.exe
```

# Adding your private key into Remote Signer

To add your private key, you can use the AddPrivateKey endpoint at `/keyRing/addPrivateKey` with the following payload:
```json
{
  "EncryptedPrivateKey": "-----BEGIN PGP PRIVATE KEY BLOCK-----\n\n (...) -----END PGP PRIVATE KEY BLOCK-----\n",
  "SaveToDisk": true,
  "Password": "12344321"
}
```

The `Password` field is optional, if provided, it will store the password along the key and auto-unlock when open. If not, you should call `UnlockKey` to be able to use the key.
