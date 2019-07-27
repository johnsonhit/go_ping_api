.PHONY: build clean

export GO111MODULE=on
export GOPROXY=https://goproxy.io

all: build

build: tidy
	go build -o bin/pingapi main.go

tidy:
	go mod tidy

clean:
	rm bin -Rf