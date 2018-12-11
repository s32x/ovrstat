deps:
	-rm -rf vendor
	-rm -f Gopkg.lock
	dep ensure
test:
	go clean
	go test ./...
clean:
	-rm -rf bin
	packr clean
run:
	make clean
	packr build -o bin/server
	./bin/server
install:
	make clean
	make deps
	packr install
deploy:
	make deps
	make test
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 packr2 build -o bin/server
	heroku container:login
	heroku container:push web -a ovrstat
	heroku container:release web -a ovrstat
	make clean