#!/usr/bin/env bash

cat << EOF > bind.go
// +build !linux,amd64
// +build !darwin,amd64
// +build !windows,amd64

// Placeholder file for assets binding
package main

func Asset(name string) ([]byte, error) {
	return []byte{}, nil
}

func AssetDir(name string) ([]string, error) {
	return []string{}, nil
}

func RestoreAssets(dir, name string) error {
	return nil
}
EOF

astilectron-bundler -d -w -l
