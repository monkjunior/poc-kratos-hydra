# Proof of concept for integration between kratos and hydra

In this example, all flows are server-side rendering types. Other types will be implemented later if needed.

There are 3 stages:

- Initialization and redirect to UI;

- Form rendering;

- Form submission and payload validation.

## Run this POC

Run example:
```bash
$ docker-compose \
  up --build --force-recreate
```

The UI will be served at `http://127.0.0.1:4455/`

Clean example:
```bash
$ docker-compose \
  rm -s -v -f
```

## Registration flow

[API Flow Golang example](https://www.ory.sh/kratos/docs/next/self-service/flows/user-registration/#registration-with-usernameemail-and-password-1)

## Login flow

## Logout flow

## References, libs and packages

Golang client package:

- [Kratos](https://pkg.go.dev/github.com/ory/kratos-client-go)

Protected endpoints by using OathKeeper, all requests come to these APIs must be authenticated.

- [Zero Trust With IAP proxy](https://www.ory.sh/kratos/docs/guides/zero-trust-iap-proxy-identity-access-proxy/)

## How to ...

- [fake json schemas ?](https://json-schema-faker.js.org/)
