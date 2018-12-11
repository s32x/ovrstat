clean:
	-rm -rf bin
	packr2 clean
	go clean
deps:
	make clean
	-rm -rf vendor
	-rm -f Gopkg.lock
	dep ensure
test:
	go test ./...
run:
	make clean
	packr2 build -o ./bin/server
	packr2 clean
	./bin/server
install:
	make clean
	make deps
	packr2 install
deploy:
	make clean
	make deps
	make test
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 packr2 build -o ./bin/server
	heroku container:login
	heroku container:push web -a ovrstat
	heroku container:release web -a ovrstat
	make clean