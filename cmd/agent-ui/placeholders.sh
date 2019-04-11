#!/usr/bin/env bash

cat << EOF > bind_darwin_amd64.go
// +build darwin,amd64

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
cat << EOF > bind_linux_amd64.go
// +build linux,amd64

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
cat << EOF > bind_windows_amd64.go
// +build windows,amd64

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
