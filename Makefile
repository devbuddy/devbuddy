PKGS=$(shell go list ./... | grep -vF /vendor/)

build:
	go build

test:
	@go test $(PKGS)

lint:
	@gometalinter.v1 \
		--vendor \
		--deadline=60s \
		--disable=gotype \
		--disable=gas \
		--exclude=".*should have comment or be unexported.*" \
		./...

docserve:
	@echo "Starting GoDoc server on http://0.0.0.0:6060"
	godoc -http=:6060

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
	gometalinter.v1 --install --update
	dep ensure
