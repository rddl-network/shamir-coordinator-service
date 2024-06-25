# Example TLS Certificates

The script in this folder generates a certificate authority as well as a server certificate used by the [Shamir Shareholder Services](https://github.com/rddl-network/shamir-shareholder-service) and a client certificate used by the Shamir Coordinator Service to authenticate against the shareholders.

## Usage

All you need to do is to call the `certs.sh` script from this folder:
```
$ ./certs.sh
-----
Certificate request self-signature ok
subject=CN = localhost
-----
Certificate request self-signature ok
subject=CN = localhost
total 24
-rw-r--r-- 1 julian users 489 May  3 09:18 ca.crt
-rw------- 1 julian users 119 May  3 09:18 ca.key
-rw-r--r-- 1 julian users 562 May  3 09:18 client.crt
-rw------- 1 julian users 119 May  3 09:18 client.key
-rw-r--r-- 1 julian users 562 May  3 09:18 server.crt
-rw------- 1 julian users 119 May  3 09:18 server.key
```

After that, point the `certs-path` in our `app.toml` for the coordinator as well as for the shareholder services to this directory.
