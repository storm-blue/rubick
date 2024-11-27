#!/bin/bash

# mac m1/m2
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o achivies/mac/arm64/rubick-darwin-arm64 .

# mac intel
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o achivies/mac/intel/rubick-darwin-amd64 .

# linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o achivies/linux/amd64/rubick-linux-amd64 .

# windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o achivies/windows/amd64/rubick-windows-amd64.exe .
