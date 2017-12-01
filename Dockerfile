##FROM golang:1.9.2-alpine

FROM ataraev/golang-alpine-git:latest

WORKDIR /go/src/locustslave

COPY . .

RUN go get github.com/golang/dep/cmd/dep && dep ensure && go build -tags 'gomq' -o bin/locustslave cmd/*

##./bin/locustslave -api_host=https://api-staging.picatic.com -task=country --master-host=127.0.0.1 --master-port=5557
#docker run cleitonmarx/locustslave ./bin/locustslave -api_host=https://api-staging.picatic.com -task=country -event_id=74448 -ticket_price_id=71362 --master-host=127.0.0.1 --master-port=5557
