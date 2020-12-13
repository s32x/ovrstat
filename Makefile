deps:
	-rm -rf ./vendor go.mod go.sum
	go mod init
	go mod tidy
	go mod vendor
	
test:
	go test ./...