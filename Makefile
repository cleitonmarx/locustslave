.PHONY: all build run

deps:
	go get -u github.com/Masterminds/glide
	glide install

build: deps
	go build -tags 'zeromq' -o bin/locustslave 

run: build
	./bin/locustslave -api_host=127.0.0.1 -event_id=123 -ticket_price_id=1 -mode=process --master-host=127.0.0.1 --master-port=5557 --rpc=zeromq
