#!/bin/bash

BUILD_LINUX_ARCH="arm arm64 386 amd64"
BUILD_OTHER_ARCH="386 amd64"
BUILD_OS="windows freebsd darwin openbsd"

TAG=`git describe --exact-match --tags HEAD`
if [[ $? -eq 0 ]];
then
  set -e
  # ----------------------------------- #
  echo "Releasing for tag ${TAG}"
  ORIGINAL_FOLDER="`pwd`"
  echo "I'm in `pwd`"
  mkdir -p zips

  # ----------------------------------- #
  echo "Installing GOX"
  go get github.com/mitchellh/gox

  # ----------------------------------- #
  echo "Building Projects"
  cd cmd
  for i in *
  do
    if [[ "${i}" != "agent-ui" ]]
    then
        echo "Building $i"
        cd $i
        echo "Running go get for linux"
        go get -v
        for os in $BUILD_OS
        do
        echo "Running go get for $os"
        GOOS=$os go get -v
        done
        mkdir -p out
        echo "Building for Linux / ${BUILD_LINUX_ARCH}"
        gox -output "out/remoteSigner-{{.Dir}}-{{.OS}}-{{.Arch}}" -arch="${BUILD_LINUX_ARCH}" -os="linux"
        echo "Building for Others / ${BUILD_OTHER_ARCH}"
        gox -output "out/remoteSigner-{{.Dir}}-{{.OS}}-{{.Arch}}" -arch="${BUILD_OTHER_ARCH}" -os="${BUILD_OS}"
        echo "Compressing builds"
        cd out
        for o in *
        do
          echo "Zipping ${o}.zip"
          zip -r "../../../zips/${o}.zip" "${o}"
        done
        cd ..
        cd ..
    fi
  done
  # ----------------------------------- #
  echo "Bundling Agent-UI"
  cd cmd/agent-ui
  ./bundleit.sh
  zip -r ../../zips/AgentUI.app.zip output/darwin-amd64/AgentUI.app
  zip -r ../../zips/AgentUI-windows-amd64.zip output/windows-amd64/AgentUI.exe
  zip -r ../../zips/AgentUI-linux-amd64.zip output/linux-amd64/AgentUI
  cd ../..
  # ----------------------------------- #
  echo "Zip Files: "
  ls -la zips
else
  echo "No tags for current commit. Skipping releases."
fi
