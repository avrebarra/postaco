PROJECTNAME := postaco

# -------------------------------------------------

## watch-web: watch webapp build
watch-web:
	parcel serve webapp/index.html -d tmp/.devweb

## setup: setup project
setup:
	GO111MODULE=off go get github.com/valyala/quicktemplate/qtc
	GO111MODULE=off go get github.com/cosmtrek/air
	GO111MODULE=off go get github.com/go-bindata/go-bindata/...
	go mod tidy
	yarn

# -------------------------------------------------

## test: test project
test:
	go test ./... -coverprofile cp.out && go tool cover -func=cp.out

## coverage: get project coverage
coverage:
	go test ./... -coverprofile cp.out && go tool cover -html=cp.out

## watch: watch golang development build
watch:
	air -c .air/development.air.toml

## prebuild: build binary
build:
	go generate ./...
	go build -o dist/$(PROJECTNAME) main.go

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run with parameter options: "
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
