#!/bin/bash

go-bindata bundle
sed -i 's/package main/package agent/g' bindata.go
