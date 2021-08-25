# Proof of concept for integration between kratos and hydra

Golang client package:

- [Kratos](https://pkg.go.dev/github.com/ory/kratos-client-go)

Protected endpoints by using OathKeeper, all requests come to these APIs must be authenticated.

- [Zero Trust With IAP proxy](https://www.ory.sh/kratos/docs/guides/zero-trust-iap-proxy-identity-access-proxy/)

Run example:
```bash
$ docker-compose \
  up --build --force-recreate
```

Clean example:
```bash
$ docker-compose \
  rm -s -v -f
```