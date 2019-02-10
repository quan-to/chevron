#!/bin/bash
BUILD_LINUX_ARCH="arm arm64 386 amd64"
BUILD_OTHER_ARCH="386 amd64"
BUILD_OS="windows freebsd darwin openbsd"

TAG=`git describe --exact-match --tags HEAD`


if [ $? -eq 0 ];
then
  set -e
  echo "Releasing for tag ${TAG}"
  ORIGINAL_FOLDER="`pwd`"
  echo "I'm in `pwd`"
  mkdir -p zips

  cd cmd
  for i in *
  do
    echo "Building $i"
    cd $i
    go get -v
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
  done
  cd ..
  echo "Zip Files: "
  ls -la zips
else
  echo "No tags for current commit. Skipping releases."
fi
