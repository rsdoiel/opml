#!/bin/bash
#
# Make releases for Linux/amd64, Linux/ARM6 and Linux/ARM7 (Raspberry Pi), Windows, and Mac OX X (darwin)
#
PROJECT=opml
VERSION=$(grep 'Version =' opml.go | cut -d\"  -f 2)
RELEASE_NAME=$PROJECT-$VERSION
echo "Preparing $RELEASE_NAME"
for PROGNAME in opmlsort opmlcat; do
  env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
  env GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
  env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
  env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
  env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/$PROGNAME.exe cmds/$PROGNAME/$PROGNAME.go
done
if [ ! -d dist ]; then
    echo "Release build failed to create dist/"
    exit 1
fi
cp -v README.md dist/
cp -v INSTALL.md dist/
cp -v LICENSE dist/
zip -r "$RELEASE_NAME-release.zip" dist/*
