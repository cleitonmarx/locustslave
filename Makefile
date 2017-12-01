.PHONY: all build run

deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure

build:
	go build -tags 'gomq' -o bin/locustslave cmd/*

run: build
	./bin/locustslave -api_host=127.0.0.1 -event_id=123 -ticket_price_id=1 -mode=process --master-host=127.0.0.1 --master-port=5557 --rpc=zeromq

