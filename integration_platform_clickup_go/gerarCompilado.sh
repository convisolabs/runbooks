#!/bin/bash

version="v03.00.03"
mainFileName="IntegrationPlatformClickup"
fileNameCompact=""
auxFileName=""

echo "Gerando Windows - AMD64"
auxFileName=$mainFileName".exe"
rm -rf $auxFileName
env GOOS=windows GOARCH=amd64 go build -o $auxFileName
fileNameCompact=$mainFileName"-"$version"-windowsAMD64.tar.xz"
tar -cJf $fileNameCompact $auxFileName
rm -rf $auxFileName
rm -rf releases/$auxFileName
mv $fileNameCompact releases/

echo "Gerando macOS - Intel"
auxFileName=$mainFileName
rm -rf $auxFileName
env GOOS=darwin GOARCH=amd64 go build -o $auxFileName
fileNameCompact=$mainFileName"-"$version"-macOSIntel.tar.xz"
tar -cJf $fileNameCompact $auxFileName
rm -rf $auxFileName
rm -rf releases/$auxFileName
mv $fileNameCompact releases/

echo "Gerando macOS - Arm64"
auxFileName=$mainFileName
rm -rf $auxFileName
env GOOS=darwin GOARCH=arm64 go build -o $auxFileName
fileNameCompact=$mainFileName"-"$version"-macOSARM64.tar.xz"
tar -cJf $fileNameCompact $auxFileName
rm -rf $auxFileName
rm -rf releases/$auxFileName
mv $fileNameCompact releases/

echo "Gerando Linux - AMD64"
auxFileName=$mainFileName
rm -rf $auxFileName
env GOOS=linux GOARCH=amd64 go build -o $auxFileName
fileNameCompact=$mainFileName"-"$version"-linuxAMD64.tar.xz"
tar -cJf $fileNameCompact $auxFileName
rm -rf $auxFileName
rm -rf releases/$auxFileName
mv $fileNameCompact releases/