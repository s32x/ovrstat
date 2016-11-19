#!/bin/bash

trap rm_ovrstat INT
function rm_ovrstat() {
        echo "Removing docker container..."
        sleep 1
        docker rm ovrstat

        echo "Removing docker image..."
        docker rmi ovrstat
}

echo "Outputting build tools info..."
docker version
docker info
go version

echo "Building binary for linux..."
env GOOS=linux GOARCH=amd64 go build -v github.com/sdwolfe32/ovrstat

echo "Building docker image..."
docker build --no-cache -t ovrstat .

echo "Cleaning up binary..."
rm ovrstat

echo "Running docker image..."
docker run --name ovrstat -p 7000:7000 ovrstat
