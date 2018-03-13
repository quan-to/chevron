#!/bin/bash

msbuild /p:Configuration=Release
rm -fr RemoteSigner/bin/Release/keys/
cd RemoteSigner/bin/Release
mkbundle -z --static --deps RemoteSigner.exe -L /usr/lib/mono/4.5 -o RemoteSigner
cd ../../../
