# BTCexporter
A lightweight Prometheus exporter that will output Bitcoin Balances for a list of addresses you specify. BTCexporter uses blockchain.info to fetch bitcoin balances. If you also wanting a Ethereum prometheus exporter, you can use [ETHexporter](https://github.com/hunterlong/btcexporter) which uses a geth server. You might also want to chart your ERC20 token balances for the Ethereum blockchain, checkout out [TOKENexporter](https://github.com/hunterlong/tokenexporter).

## Watch Addresses
The `addresses.txt` file holds all the addresses to fetch balances for. Use the format `name:address` on each new line. BTCexporter updates balances every 60 seconds since the bitcoin blockchain doesn't update very quickly.
```
example2:1Kr6QSydW9bFQG1mXiPNNu6WpJGmUa9i1g
example3:17A16QmavnUfCW11DAApiJxp7ARnxN5pGX
```

## Build Docker Image
Clone this repo and then follow the simple steps below!

##### Build Docker Image
`docker build -t hunterlong/btcexporter:latest .`

##### Run ethexporter
`docker run -d -p 9019:9019 hunterlong/btcexporter:latest`

## Pull from Dockerhub
Create a `addresses.txt` file with the correct format mentioned above.
```
docker run -d -v ~/btcexporter:/app \
 -p 9019:9019 \
 hunterlong/btcexporter:latest
```
The Docker image should be running with the default addresses.

## Prometheus Response
```
btc_balance{name="example2",address="1Kr6QSydW9bFQG1mXiPNNu6WpJGmUa9i1g"} 2543.029143
btc_balance{name="example3",address="17A16QmavnUfCW11DAApiJxp7ARnxN5pGX"} 7286.88533
btc_balance{name="example4",address="1DcKsGnjpD38bfj6RMxz945YwohZUTVLby"} 1271.125171
btc_balance{name="example5",address="3JjPf13Rd8g6WAyvg8yiPnrsdjJt1NP4FC"} 129.7400535
btc_balance{name="example6",address="1TjstSNNZezhTMj6m9pcGwMr1fxLhwUuH"} 0
btc_balance{name="example7",address="1NDyJtNTjmwk5xPNhjgAMu4HDHigtobu1s"} 42094.21118
btc_balance{name="example8",address="3DzSVk4veMCkNbNT9CdETeE26uWxmNbBnD"} 1243.824953
btc_balance{name="example9",address="1LV5y3NkVkmdWnF6xRCEXrAnUkRgge4KSq"} 529.1578158
btc_balance{name="example10",address="1EEqRvnS7XqMoXDcaGL7bLS3hzZi1qUZm1"} 1377.11966
```
