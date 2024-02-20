# shamir-coordinator-service
This service serves the purpose of creating and distributing shamir secret shares to a set of distributed shareholder services, and collecting them again in order to sign a transaction.
The transation will be issued on Liquid. The elemnts RPC can be configured with `app.toml` file.
To ensure secure communication it utilizes mutual TLS. It offers two routes:

- POST `/send/:recipient/:amount`
- POST `/mnemonics/:secret`

## Prerequisits
The creation and the recovery of the shares is done with the help of `https://github.com/trezor/python-shamir-mnemonic`. Therefore, the python needs to be installed and the virtual env path of the related python environment/binary needs to be defined in the configuration. 

## Execution

The service can be executed via the following go command without having it previously built:
```bash
go run cmd/shamir-coordinator/main.go
```

## Configuration
The service needs to be configured via the ```./app.toml``` file or environment variables. The defaults are
```
asset-id = 'asset-id'
certs-path = './certs/'
rpc-enc-timeout = 20
rpc-host = 'localhost'
rpc-password = 'password'
rpc-port = 18884
rpc-scheme = 'http'
rpc-user = 'user'
rpc-wallet-name = 'wallet'
service-bind = 'localhost'
service-port = 8080
shamir-shares = 3
shamir-threshold = 2
share-holder-list = 'https://localhost:8081,https://localhost:8082,https://localhost:8083'
test-mode = false
virtual-env-path = '/opt/hostedtoolcache/Python/3.10.13/x64'
```
