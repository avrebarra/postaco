PROJECTNAME := postaco

## watch-web: watch webapp build
watch-web:
	parcel serve webapp/index.html -d tmp/.devweb

## setup: setup project
setup:
	go get -u github.com/valyala/quicktemplate/qtc
	go get -u github.com/cosmtrek/air
	go get -u github.com/go-bindata/go-bindata/...
	yarn

## test: test project
test:
	go test ./... -coverprofile cp.out && go tool cover -func=cp.out

## coverage: get project coverage
coverage:
	go test ./... -coverprofile cp.out && go tool cover -html=cp.out

## watch: watch golang development build
watch:
	air -c .air/development.air.toml

## build: build binary
build: setup
	go generate .
	go build -o dist/$(PROJECTNAME) main.go

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run with parameter options: "
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
