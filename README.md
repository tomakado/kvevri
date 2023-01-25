# kvevri

Key value data store with gRPC interface.

# Usage

## On host

```
make run
```

or

```
make build && bin/kvevri
```

## With Docker

```
docker build -t kvevri:latest .
docker run -p 8080:8080 kvevri:latest
```

# Config

It's possible to configure kvevri with environment variables:

* `KVEVRI_LISTEN_ADDR` — address to listen in `<host>:<port>` format, default `8080`
* `KVEVRI_TTL` — how long keys should be preserved, default `1h`

