#!/bin/bash

# mac m1/m2
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o achivies/mac/arm64/rubick ./tools/cmd/

# mac intel
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o achivies/mac/intel/rubick ./tools/cmd/

# linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o achivies/linux/amd64/rubick ./tools/cmd/

# windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o achivies/windows/amd64/rubick.exe ./tools/cmd/
