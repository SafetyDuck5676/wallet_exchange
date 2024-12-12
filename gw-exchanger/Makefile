.PHONY: build run docker-build docker-run

build:
 GOOS=linux GOARCH=amd64 go build -o main ./cmd

run:
 ./main -c config.env

docker-build:
 docker build -t gw-exchanger .

docker-run:
 docker run --rm -p 50051:50051 --env-file config.env gw-exchanger