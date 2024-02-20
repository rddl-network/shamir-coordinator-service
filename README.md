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
ServiceBind:     "localhost",
ServicePort:     8080,
ShareHolderList: "https://localhost:8081,https://localhost:8082,https://localhost:8083",
CertsPath:       "./certs/",
RpcWalletName:   "wallet",
RpcHost:         "localhost",
RpcPort:         18884,
RpcUser:         "user",
RpcPassword:     "password",
RpcScheme:       "http",
RpcEncTimeout:   20,
AssetID:         "asset-id",
ShamirThreshold: 2,
ShamirShares:    3,
VirtualEnvPath:  "/opt/hostedtoolcache/Python/3.10.13/x64",
TestMode:        false,
```

