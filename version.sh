#!/bin/bash

cat << EOF > golem/version.go
package kitchen

const (
	version = "$(git describe --tags)"
)
EOF