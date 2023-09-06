#!/bin/bash

# build mac-intel
GOOS=darwin GOARCH=amd64 go build -o dist/intel_mac/curly ./cmd/curly/

# build mac-m1
GOOS=darwin GOARCH=amd64 go build -o dist/m1_mac/curly ./cmd/curly/

# build windoze 32
GOOS=windows GOARCH=386 go build -o dist/win32/curly ./cmd/curly/

# build windoze 64
GOOS=windows GOARCH=amd64 go build -o dist/win64/curly ./cmd/curly/