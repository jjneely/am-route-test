.PHONY: build

export GO111MODULE=on
export CGO_ENABLED=0


build:
	go build -mod=vendor -o am-route-test
