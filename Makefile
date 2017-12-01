build:
	go build

run: build
	./dad

install:
	go install

build-release:
	go build -ldflags "-s -w"
	upx dad  # brew install upx
