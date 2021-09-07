# Proof of concept for integration between kratos and hydra

## Run this POC

Start docker-compose stack
```bash
$ make start-all
```

A login and consent app will be served at `http://127.0.0.1:4455/`

Create the first Hydra client
```bash
$ make create-hydra-client
```

Start an exemplary client to perform OAuth2 Authorization Code flow.
```bash
$ make examine-authorization-code
```
Visit the client at `http://127.0.0.1:5555`

Rebuild and update UI service:
```bash
$ make update-ui
```

Clean example:
```bash
$ make clean-all
```

# My personal docs for this example

## OIDC and Hydra

[Kratos configuration](https://www.ory.sh/kratos/docs/concepts/credentials/openid-connect-oidc-oauth2/#configuration)

Discussions about this topic:

- [kratos-1145](https://github.com/ory/kratos/discussions/1145)

- [kratos-1511](https://github.com/ory/kratos/discussions/1511)

## References, libs and packages

[Cookies vs Tokens](https://dzone.com/articles/cookies-vs-tokens-the-definitive-guide)

Golang client package:

- [Kratos](https://pkg.go.dev/github.com/ory/kratos-client-go)

- [Hydra golang SDK](https://www.ory.sh/hydra/docs/sdk/go/)

- Enable OpenID connect support for `golang.org/x/oauth2` by using `github.com/coreos/go-oidc`

Protected endpoints by using OathKeeper, all requests come to these APIs must be authenticated.

- [Zero Trust With IAP proxy](https://www.ory.sh/kratos/docs/guides/zero-trust-iap-proxy-identity-access-proxy/)

## How to ...

- [fake json schemas ?](https://json-schema-faker.js.org/)
