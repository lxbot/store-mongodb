.PHONY: build

build:
	go build -buildmode=plugin -o store-mongodb.so store.go