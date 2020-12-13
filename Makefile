deps:
	-rm -rf ./vendor go.mod go.sum
	go mod init
	go mod tidy
	go mod vendor
	
test:
	go test ./...

deploy: deps test
	up prune -s production -r 2
	-up stack plan
	-up stack apply
	up deploy production