#!/bin/bash

cat << EOF > kitchen/version.go
package kitchen

const (
	version = "$(git describe --tags)"
)
EOF