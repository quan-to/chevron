Quanto Remote Signer (QRS)
====================

[![MIT License](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://tldrlegal.com/license/mit-license)

A simple Web Server to act as a GPG Creator / Signer / Verifier. This abstracts the use of the GPG and makes easy to sign / verify any GPG document using just a POST request.

Please notice that this application is *NOT inteded to ran public in the internet*. This is inteded to be a helper service to your application be able to sign / verify data (same as local gpg in the system). Because of that, it only listens for localhost.

It is based off in [racerxdl](https://github.com/racerxdl) [AppServer](https://github.com/racerxdl/AppServer) to provide a REST web server. It is authorized by the owner to be licensed at MIT here (check the Owner Signature at the commit that adds this README as the proof)

TODO
====

*   Increment this document with models and enums
*   Encrypt / Decrypt calls

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

For signing data, the first thing is to make sure the key is loaded and decrypted. Then encode the data you want to sign in Base64. Let's take for example the following text to sign:

```
Terrible connections lead to terrible users
```

We can base64 encode (for example in shellscript):
```bash
echo "Terrible Connections lead to Terrible Users" | base64
VGVycmlibGUgQ29ubmVjdGlvbnMgbGVhZCB0byBUZXJyaWJsZSBVc2Vycwo=
```

And make a POST request to `/remoteSigner/gpg/sign` with the following JSON payload:
```json
{
  "Base64Data": "VGVycmlibGUgQ29ubmVjdGlvbnMgbGVhZCB0byBUZXJyaWJsZSBVc2Vycwo=",
  "FingerPrint": "D7362B4CC546DB11"
}
```
Leading to a Ascii Armored GPG Detached Signature:
```
-----BEGIN PGP SIGNATURE-----
Version: BCPG C# v1.8.1.0

iQIcBAABCgAGBQJaaldaAAoJENc2K0zFRtsRSFoP+gIf8iKOX0+jOJzXis8yGdb6
yYiDq2iKKhLKrMR110F5gWOkbizLVut47XuR9zj6RBWvE2yd2YUJS/bom2A1F+4n
UZMdwvk0frpsoJPbwu5wjHKUuWTar3H8dgd5Wdx3BaNlr/j2JvTOq6/MGBzBJQ9y
39vve5KQnMzimKcKSUga2RY/5dUpGGJFvWKB1KNoh/jGBG7LxqA2zVh+o1b8cWcE
BW9Qk1fq5j9q+N68LJ9SGuHXALrxVXeR6Z9GGQjcPIwZbTYhPt1IKgr20pvvCc/6
GaXRSEkcacike+zNjLgMewL4k3/wRMceVEq+FKT1RtGciEr6dd+WsAbPA7DzoOmn
XtrkCzCFNay9F7AHM1eopYP5qROEOCIBi+vXPLLAsZSua8cPd8DUxmNDDAc9Cq+F
DCD/v3FZXoqmTmNaJ5QJi6hW8lvUMK9KDzN51CmlOoICyu8BInPZIf0VO4NLlD/d
v9CR+mHaHWmYcSEpoI4I+YpNB9JmxwIFd8+KluqkQ5tzixBeVx5O0gO2yuStG0Le
mZTOf886+cl49KwOBzZZS37lx6VCxgMTI3dgUL+4r4L1em4QmOjEoNu0fCbigkoP
KZ3OYGRhhp3dsmIIBuOWC+6PvbzNKApaph4wH4ysLZDP7/DEtPLC5msuEYW5QcJt
TPTWS9moWewsplQazivx
=fEEJ
-----END PGP SIGNATURE-----
```

You might also want to get directly into *Quanto Signature* format (`FingerPrint_Hash_Signature`). For doing so you can send the same payload to `/remoteSigner/gpg/signQuanto` leading to:

```
D7362B4CC546DB11_SHA512_iQIcBAABCgAGBQJaale3AAoJENc2K0zFRtsRe3IP/jVE19IeT9fWXl1wbSfZ4VLRY8HePogfyGMVELrkqoRQjUwQB3s2cBio/uAZNNzyvYGqkdFVeeSO83GRAsobts8Q94Q//jAJxeYDy6qAzs6JbzOYAf1b8KWhjzosQDnvmqlvyH+95IoxTEXcDK/WFox5XrZGqRda3rlv+9CywzYreAiFHnSuF5LFJ0K+KkPCMjEJ8EgRZQ/WN0gcDNcabgI85ncpJ7gQ6rSzOmvK3tDd3oNyFFfzYGNaGWThQsYKLOwZA3MSri95y86CcBW9SkaLdqT9LRSGW+pEjXYXax4WU13+YlrUa5axT87sHZs0awORKkHZ2Wik082cFN7M903qd1+fUkKNaz3nG1rwbUAp5KiKabQUcvOhz+guXYnZlqeL0IWRBvagAnBDjWLd3O4X9RIhhl4RjqeiHzgzpx3hBbUQRxyDopdnQMGuqIH9PaffJzFUnzqChTgUhnxntYqkISPsYy5DMzzvlLIIESvIgCHOgQu9kxj/upTE6OoDjWMXjBJn3ytpwf2xjsdtKEn0QRe0PV1uCa+P+z5Qg41ZvP0krhxomr1wmNmNvDkPL/uIYD06fN6bWwnuMQf5nI5DS/X4ysp1AN5EEesrwjh9ygBGxpFa4+jlwuJYa7b2HdmesQG3JVzMhkHtNCCMIo7GAHKCc8vhG08eQ0FtAdPr=rV18
```

#### Verifing Signatures

The process to verify a signature is similar to signing data. First you need to make sure the public key is loaded into the server (if the private key is, the public is as well) or it is available in the specified SKS Server and then encode the test data in Base64 format and POST the following payload to `/remoteSigner/gpg/verifySignature`

```json
{
  "Base64Data": "VGVycmlibGUgQ29ubmVjdGlvbnMgbGVhZCB0byBUZXJyaWJsZSBVc2Vycwo=",
  "signature":"-----BEGIN PGP SIGNATURE-----\nVersion: BCPG C# v1.8.1.0\n\niQIcBAABCgAGBQJaalMSAAoJENc2K0zFRtsRuxEQAKPz46GpBYvZY99dqfylo/ux\nOcbFx0U/jWnXACEsz4KbfIaKqTNQLNOApt7vC+PTeK18Bx3i6lLDq5s0T56ZZS07\n+12/8qWfq08LOTANtrMetmOP1znaJpzmWzpxCp8t/pTowt1RZJUfGC5zdxoxMLB8\nu3siSbbqSxPlOYx7yfNnJqE7KagHn83WdZYIQTFBYZESqfEhjmazERui+g3YKF74\nUlr8Ey0kIFptVa/DdsIQwgCMjDalWB6zdX8xLiLqH0pRAOSmiMcU7cX8vWSAlC5S\nl7uoam6azzyc/kAYBaBd+3/YORDu6vlnNvHtj6D1cvff0ahinnKHb2gxqN+cBUbL\nNlpBSRKwd5i/O504BjGOplp31rCbza+0vs0lvbQ2/ZMH1HaspuO3I8jdZQpaMxR4\n6GWf9/+clg8SxkNgbaKj8pRHnzvrjEOaEfYNXxdjMV4LgXX6pjZtZi47AfHqXPWc\n/pw7YSG6yJLF4n+Egky/thvoVQHR3GM10VXvUjYVTRdEaxO2P9MpWYLD8Fqa6rMo\nzVNXtCMKtej5Y5qQHMCFcKafG1J0TXPfsqnYCSsG2PpdDsZhz5r6eKgu5LZqN5yO\ntRylJnBDJuE8yweNYBDYWWAMoo1ApziRUItNI5el1I0DC+vwaM4Vxurywjxhcy4q\nHTlBXAdXKpVLINxkRSKy\n=YK2F\n-----END PGP SIGNATURE-----" 
}
```

In case of success, it will return an `OK`. In case of failure it will return an `ErrorObject` with an `errorCode` field `INVALID_SIGNATURE`.

You might also want to verify in `Quanto Signature` format. To do so, you can send the same payload with the `signature` field containing the `Quanto Signature` you want to verify to `/remoteSigner/gpg/verifySignatureQuanto`.

```json
{
  "Base64Data": "VGVycmlibGUgQ29ubmVjdGlvbnMgbGVhZCB0byBUZXJyaWJsZSBVc2Vycwo=",
  "signature":"D7362B4CC546DB11_SHA512_iQIcBAABCgAGBQJaale3AAoJENc2K0zFRtsRe3IP/jVE19IeT9fWXl1wbSfZ4VLRY8HePogfyGMVELrkqoRQjUwQB3s2cBio/uAZNNzyvYGqkdFVeeSO83GRAsobts8Q94Q//jAJxeYDy6qAzs6JbzOYAf1b8KWhjzosQDnvmqlvyH+95IoxTEXcDK/WFox5XrZGqRda3rlv+9CywzYreAiFHnSuF5LFJ0K+KkPCMjEJ8EgRZQ/WN0gcDNcabgI85ncpJ7gQ6rSzOmvK3tDd3oNyFFfzYGNaGWThQsYKLOwZA3MSri95y86CcBW9SkaLdqT9LRSGW+pEjXYXax4WU13+YlrUa5axT87sHZs0awORKkHZ2Wik082cFN7M903qd1+fUkKNaz3nG1rwbUAp5KiKabQUcvOhz+guXYnZlqeL0IWRBvagAnBDjWLd3O4X9RIhhl4RjqeiHzgzpx3hBbUQRxyDopdnQMGuqIH9PaffJzFUnzqChTgUhnxntYqkISPsYy5DMzzvlLIIESvIgCHOgQu9kxj/upTE6OoDjWMXjBJn3ytpwf2xjsdtKEn0QRe0PV1uCa+P+z5Qg41ZvP0krhxomr1wmNmNvDkPL/uIYD06fN6bWwnuMQf5nI5DS/X4ysp1AN5EEesrwjh9ygBGxpFa4+jlwuJYa7b2HdmesQG3JVzMhkHtNCCMIo7GAHKCc8vhG08eQ0FtAdPr=rV18" 
}
```

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

Environment Variables
=====================

These are the Environment Variables that you can set to manage the webserver:

*   `SYSLOG_IP` => IP of the Syslog Server to send Console Messages _(defaults to '127.0.0.1')_ *Does not apply for Windows*
*   `SYSLOG_FACILITY` => Facility of the Syslog to use. _(defaults to 'LOG_USER')_
*   `PRIVATE_KEY_FOLDER` => Folder to load / store encrypted private keys. _(defaults to './keys')_
*   `SKS_SERVER` => SKS Server to fetch / put public keys. _(defaults to 'http://pgp.mit.edu/')_
*   `MAX_KEYRING_CACHE_SIZE` => Maximum Number of Public Keys to cache (does not include Private Keys derived Public Keys). _(defaults to 1000)_


Building
========

# For Windows
Just use the newest version of Visual Studio and hit build. You should have a executable in folder `bin/Debug` or `bin/Release`.

# For Linux
You will need `mono` and `nuget`. On `Ubuntu 16.04` you can run:
```bash
sudo apt install mono-complete nuget
git clone https://github.com/quan-to/remote-signer
cd remote-signer
nuget restore
msbuild /p:Configuration=Release
mkdir binaries
cp RemoteSigner/bin/Release/* binaries
mkdir binaries/keys
cd binaries
```
Aditionally, you can compile an executable that does not need mono to be ran by using mkbundle:

```bash
mkbundle -z --static --deps RemoteSigner.exe -L /usr/lib/mono/4.5 -o RemoteSigner
```

This will generate a static binary called `RemoteSigner`.

# For Mac

You will need `mono` and `nuget`. With [Brew](https://brew.sh/) you can run:

```bash
brew install mono nuget
git clone https://github.com/quan-to/remote-signer
cd remote-signer
nuget restore
xbuild /p:Configuration=Release
mkdir -p binaries/keys
cp RemoteSigner/bin/Release/* binaries
cd binaries
open ./
```
