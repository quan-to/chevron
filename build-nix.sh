#!/bin/bash

nuget restore
msbuild /p:Configuration=Release /t:Clean,Build
rm -fr RemoteSigner/bin/Release/keys/
cd RemoteSigner/bin/Release
cd ../../../
