#!/bin/sh
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export COLOR='\033[0;32m' # Green
export NC='\033[0m' # No Color

printf "${COLOR}DELETING OLD DOCKER CONTAINER${NC}\n"
docker kill ovrstat
docker rm ovrstat

printf "${COLOR}RETRIEVING DEPENDENCIES${NC}\n"
go get

printf "${COLOR}BUILDING FRESH BINARY${NC}\n"
go build -o bin/ovrstat

printf "${COLOR}BUILDING DOCKER CONTAINER${NC}\n"
docker build --no-cache -t sdwolfe32/ovrstat .

printf "${COLOR}CLEAING UP BINARY${NC}\n"
rm -r bin/
