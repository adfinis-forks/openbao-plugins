# Consul secrets engine (API)

This is the API documentation for the OpenBao Consul secrets engine. For general
information about the usage and operation of the Consul secrets engine, please see the
[documentation](../readme.md).

This documentation assumes the Consul secrets engine is mounted at the `/consul` path
in OpenBao. Since it is possible to mount secrets engines at any location, please
update your API calls accordingly.

## Configure access

This endpoint configures the access information for Consul. This access
information is used so that OpenBao can communicate with Consul and generate
Consul tokens.

| Method | Path                    |
| :----- | :---------------------- |
| `POST` | `/consul/config/access` |

### Parameters

- `address` (`string: "127.0.0.1:8500"`) - Consul server address

- `token` (`string`) - Consul management token. If missing ACL bootstrapping
  will be attempted.

- `scheme` (`string: "http"`) - URI scheme ("http" or "https")

- `ca_cert` (`string`) - CA certificate to use when verifying Consul server
  certificate, must be x509 PEM encoded.

- `client_cert` (`string`) - Client certificate used for Consul's TLS
  communication, must be x509 PEM encoded and if this is set you need to also
  set `client_key`.

- `client_key` (`string`) - Client key used for Consul's TLS communication, must
  be x509 PEM encoded and if this is set you need to also set `client_cert`.


### Sample payload

```json
{
  "address": "consul.example.com:8500",
  "token": "ffffffff-ffff-ffff-ffff-ffffffffffff",
  "scheme": "https"
}
```

### Sample request

```shell-session
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    http://127.0.0.1:8200/v1/consul/config/access
```

## Read access configuration

This endpoint queries for information about the Consul connection.

| Method | Path                    |
| :----- | :---------------------- |
| `GET`  | `/consul/config/access` |

### Sample request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/consul/config/access
```

### Sample response

```json
{
  "data": {
    "address": "consul.example.com:8500",
    "scheme": "https"
  }
}
```

## Create/Update role

This endpoint creates or updates the Consul role definition in OpenBao. If the
role does not exist, it will be created. If the role already exists, it will
receive updated attributes.

| Method | Path                  |
| :----- | :-------------------- |
| `POST` | `/consul/roles/:name` |

### Parameters

- `name` (`string: <required>`) – Specifies the name of the role to
  create/update. This is part of the request URL.

- `consul_policies` (`string`) - Comma-separated list of policies to attach to the
  token. Either `consul_policies` or `consul_roles` are required for Consul 1.5
  and above, or just `consul_policies` if using Consul 1.4.

- `consul_roles` (`string`) - Comma-separated list of Consul roles to attach to
  the token. Either `consul_policies` or `consul_roles` are required for Consul
  1.5 and above.

- `consul_namespace` (`string: "default"`) - Indicates which namespace that the
  token will be created within. Available in Consul 1.7 and above.

- `node_identities` (`string`) - Comma-separated list of Node Identities to
  attach to the token. Available in Consul 1.8.1 or above.

- `partition` (`string: "default"`) - Indicates which admin partition that the
  token will be created within. Available in Consul 1.11 and above.

- `service_identities` (`string`) - Comma-separated list of Service Identities
  to attach to the token. Available in Consul 1.5 or above.

- `local` (`bool: false`) - Indicates that the token should not be replicated
  globally and instead be local to the current datacenter.

- `ttl` (`duration: 1h`) - Token time-to-live

- `max_ttl` (`duration: 24h`) - Maximum token time-to-live

### Sample payload

To create a client token with a custom policy:

```json
{
  "consul_policies": "policy1,policy2",
  "ttl": "10m",
  "max_ttl": "1h"
}
```

### Sample request

```shell-session
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    http://127.0.0.1:8200/v1/consul/roles/worker
```

## Read role

This endpoint queries for information about a Consul role with the given name.
If no role exists with that name, a 404 is returned.

| Method | Path                  |
| :----- | :-------------------- |
| `GET`  | `/consul/roles/:name` |

### Parameters

- `name` (`string: <required>`) – Specifies the name of the role to query. This
  is part of the request URL.

### Sample request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/consul/roles/worker
```

### Sample response

```json
{
  "data": {
    "consul_namespace": "",
    "consul_policies": ["policy1", "policy2"],
    "local": false,
    "max_ttl": 3600,
    "partition": "",
    "ttl": 600
  }
}
```

## List roles

This endpoint lists all existing roles in the secrets engine.

| Method | Path                      |
| :----- | :------------------------ |
| `LIST` | `/consul/roles`           |
| `GET`  | `/consul/roles?list=true` |

### Sample request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/consul/roles
```

### Sample response

```json
{
  "data": {
    "keys": ["worker"]
  }
}
```

## Delete role

This endpoint deletes a Consul role with the given name. Even if the role does
not exist, this endpoint will still return a successful response.

| Method   | Path                  |
| :------- | :-------------------- |
| `DELETE` | `/consul/roles/:name` |

### Parameters

- `name` (`string: <required>`) – Specifies the name of the role to delete. This
  is part of the request URL.

### Sample request

```shell-session
$ curl \
    --request DELETE \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/consul/roles/example-role
```

## Generate credential

This endpoint generates a dynamic Consul token based on the given role
definition.

| Method | Path                  |
| :----- | :-------------------- |
| `GET`  | `/consul/creds/:name` |

### Parameters

- `name` (`string: <required>`) – Specifies the name of an existing role against
  which to create this Consul token. This is part of the request URL.

### Sample request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/consul/creds/example
```

### Sample response

```json
{
  "data": {
    "accessor": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "consul_namespace": "",
    "local": false,
    "partition": "",
    "token": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
  }
}
```
