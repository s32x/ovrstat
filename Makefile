clean:
	-rm -rf bin
	go clean
deps:
	make clean
	-rm -rf vendor
	-rm -r glide.yaml
	-rm -f glide.lock
	glide init --non-interactive
	glide install
test:
	go test ./...
run:
	make clean
	go build -o ./bin/server
	go clean
	./bin/server
install:
	make clean
	make deps
	go install