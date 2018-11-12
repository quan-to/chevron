#!/bin/bash

rm -fr ~/.nuget/packages
rm -fr packages
nuget update -self
nuget locals all -clear
nuget restore
msbuild /p:Configuration=Debug /t:Clean,Build
rm -fr RemoteSigner/bin/Debug/keys/
cd RemoteSigner/bin/Debug
mkbundle -z --static --deps RemoteSigner.exe -L /usr/lib/mono/4.5 -o RemoteSigner
cd ../../../
