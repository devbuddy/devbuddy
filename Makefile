build:
	go build

test:
	@go test `go list ./... | grep -vF /vendor/`

lint:
	@gometalinter.v2 \
		--vendor \
		--deadline=120s \
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

releases:
	GOARCH=amd64 GOOS=darwin CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s -w " -o dist/dad-darwin-amd64
	shasum -a 256 dist/dad-darwin-amd64 > dist/dad-darwin-amd64.sha256

	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s -w " -o dist/dad-linux-amd64
	shasum -a 256 dist/dad-linux-amd64 > dist/dad-linux-amd64.sha256

	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s -w " -o dist/dad-windows-amd64.exe
	shasum -a 256 dist/dad-windows-amd64.exe > dist/dad-windows-amd64.exe.sha256

devup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u gopkg.in/alecthomas/gometalinter.v2
	gometalinter.v2 --install --update
	dep ensure
