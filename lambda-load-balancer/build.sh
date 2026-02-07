#!/bin/bash

echo "Build the binary"
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

echo "Create a zip file"
zip deployment.zip bootstrap

echo "Cleaning up"
rm bootstrap