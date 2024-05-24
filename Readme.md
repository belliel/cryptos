# Exercise

```
App uses Kraken API to show close price of asset pairs (like BTC/UTC)

/api/v1/ltp
    - pairs=BTC/USD
        - default pairs = BTC/USD,BTC/CHF,BTC/EUR
```

How to run app
```
cryptos [OPTIONS]

Application Options:
      --http-port=                                                http binding port (default: 8080) [$HTTP_PORT]
      --crypto-service-target-last-traded-prices-pairs=
      --crypto-service-target-last-traded-prices-worker-interval=
      --crypto-service-cache-ttl=
      --log-level=[debug|error|warn|info]

Help Options:
  -h, --help                                                      Show this help message
```

How to build:
```sh
make build
```

How to docker bundle:
```sh
make docker
```

How to run tests
```sh
make test
```