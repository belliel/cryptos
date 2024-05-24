docker:
	docker build -t cryptos:latest ./

build:
	go build -o build/cryptos cmd/app/main.go