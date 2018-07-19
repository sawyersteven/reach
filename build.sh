#!/bin/bash

echo Building darwin.386
GOOS=darwin GOARCH=386 go build -o ./releases/darwin.386/reach

echo Building linux.386
GOOS=linux GOARCH=386 go build -o ./releases/linux.386/reach

echo Building windows.386
GOOS=windows GOARCH=386 go build -o ./releases/windows.386/reach.exe

