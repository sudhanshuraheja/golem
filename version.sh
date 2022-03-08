#!/bin/bash

cat << EOF > golem/version.go
package golem

const (
	version = "$(git describe --tags)"
)
EOF