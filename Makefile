clean:
	-rm -rf bin
	packr2 clean
deps:
	make clean
	-rm -rf vendor
	-rm -r glide.yaml
	-rm -f glide.lock
	glide cache-clear
	glide init --non-interactive
	glide install
test:
	go test ./...
run:
	make clean
	packr2 build -ldflags="-s -w" -o ./bin/server
	./bin/server
install:
	make clean
	make deps
	packr2 install