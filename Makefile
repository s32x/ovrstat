init:
	-rm -rf ./vendor go.mod go.sum
	GO111MODULE=on go mod init

deps:
	-rm -rf ./vendor go.sum
	GO111MODULE=on go mod vendor
	
test:
	go test ./...

deploy: deps test
	up prune -s production -r 2
	-up stack plan
	-up stack apply
	up deploy production