########################################################
-include Makefile.options.local
port?=
GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
########################################################
serve:
	PORT=$(port) go run main.go

build:
	$(GO_BUILD_ENV)go build -o bin/shakesearch

deploy:
	git push heroku master

test:
	go test ./...
########################################################
