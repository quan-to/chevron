#!/bin/bash

rm -fr ~/.nuget/packages
rm -fr packages
nuget update -self
nuget locals all -clear
nuget restore
msbuild /p:Configuration=Release /t:Clean,Build
rm -fr RemoteSigner/bin/Release/keys/
cd RemoteSigner/bin/Release
mkbundle -z --static --deps RemoteSigner.exe -L /usr/lib/mono/4.5 -o RemoteSigner
cd ../../../
