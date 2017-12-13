PKGS=$(shell go list ./... | grep -vF /vendor/)

build:
	go build

test:
	@go test $(PKGS)

lint:
	@gometalinter.v1 --vendor \
		--disable=gotype \
		--disable=gas \
		--exclude=".*should have comment or be unexported.*" \
		./...

run: build
	./dad

install:
	go install

build-release:
	go build -ldflags "-s -w"
	upx dad  # brew install upx

devup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u gopkg.in/alecthomas/gometalinter.v1
	dep ensure
