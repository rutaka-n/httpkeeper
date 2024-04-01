# HTTP Keeper
The HTTP Keeper is a simple http reverse proxy server that proxy http-requests to an upstream server.
It allows to:
- set basic auth user:password for upstream
- use access-tokens to authenticate clients requests
- set the rate limit for upstream
- set invalidated tokens

## Installation

Download the archive from release page, extract the executable file, and add into your path.

## Configuration

You can find the example of configuration file in `share/config.json`
**Note:** Please change the `server.secret`.

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
* TLS
* control script
* install.sh
* Load balancing
* healthchecks

