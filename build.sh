#!/bin/bash

# build mac-intel
GOOS=darwin GOARCH=amd64 go build -o dist/intel_mac/curly ./cmd/curly/

# build mac-m1
GOOS=darwin GOARCH=arm64 go build -o dist/m1_mac/curly ./cmd/curly/

# build windoze 32
GOOS=windows GOARCH=386 go build -o dist/win32/curly.exe ./cmd/curly/

# build windoze 64
GOOS=windows GOARCH=amd64 go build -o dist/win64/curly.exe ./cmd/curly/

# zip for windoze
cd ./dist

zip -r curly_win32.zip ./win32/curly.exe

zip -r curly_win64.zip ./win64/curly.exe