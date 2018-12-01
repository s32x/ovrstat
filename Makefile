deps:
	-rm Gopkg.toml
	-rm Gopkg.lock
	-rm -r vendor
	dep init
test:
	go clean
	go test ./...
run:
	go run main.go
deploy:
	make deps
	make test
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/server
	heroku container:login
	heroku container:push web -a ovrstat
	heroku container:release web -a ovrstat
	rm -r bin