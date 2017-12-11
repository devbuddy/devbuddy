PKGS=$(shell go list ./... | grep -vF /vendor/)

build:
	go build

test:
	go test $(PKGS)

lint:
	go vet $(PKGS)
	golint $(PKGS)

run: build
	./dad

install:
	go install

build-release:
	go build -ldflags "-s -w"
	upx dad  # brew install upx
