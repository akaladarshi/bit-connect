# bit-connect

## Description
This is a simple project to do a handshake with a bitcoin node using Go.

# Features
- Simple project to connect to a bitcoin node.
- Connect to a local bitcoin node at default address `0.0.0.0:18444`.
- To connect to a different address, use the flag `--node-address`
- By default, the project will use 76001 as the protocol version.
- By default, the project will connect to the `regtest` network.
- Only support connection handshake then disconnects.
- Only IPv4 addresses are supported.

## Prerequisites
- Go version 1.22.1
- logging `github.com/rs/zerolog v1.32.0`
- commands `github.com/spf13/cobra v1.8.0`

## How to run

### Run a bitcoin node
```shell
git clone https://github.com/btcsuite/btcd

make install 

btcd --regtest
```

### Run the project
```shell
make connect 

or 

make build

./bin/bit-connect connect --node-address = "0.0.0.0:18444"
```

