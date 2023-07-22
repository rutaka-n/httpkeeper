# HTTP Keeper
The HTTP Keeper is a simple http reverse proxy server that proxy http-requests to an upstream server.
It allows to:
- set basic auth user:password for upstream
- use access-tokens to authenticate clients requests
- set the rate limit for upstream
- set invalidated tokens

## Instalation

## Configuration

## Usage

### Generate Secret

```
$ httpkeeper secret N
```
Where `N` is a length of the secret.
After generating, you need to set the value to the config file.

### Generate Token

```
$ httpkeeper token -config path/to/config.json -client CLIENT_NAME -expiresAt "2023-12-31 23:59:59"
```
The command generate JWT-token that you can copy to configuration of the client. Client should set the value in HTTP Authorization header:
```
Authorization: Bearer <token>
```

### Run proxy

```
$ httpkeeper proxy -config path/to/config.json [-addr addr:port] [-log path/to/httpkeeper.log]
```

## TODO
* set/check list of services in the token
* TCP timeouts
* TLS
* control script
* install.sh
* Load balancing
* healthchecks

