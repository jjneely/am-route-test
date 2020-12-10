.PHONY: build

export GO111MODULE=on
export CGO_ENABLED=0


build:
	go build -o am-route-test

clean:
	rm am-route-test
