# Consul secrets engine

Name: `Consul`

Consul is a service networking platform featuring service discovery, service
mesh, API gateway and configuration management. The Consul secrets engine for
OpenBao generates [Consul](https://developer.hashicorp.com/consul) ACL tokens
dynamically based on pre-existing Consul ACL policies.

This page will show a quick start for this secrets engine. For detailed
documentation on every path, use `bao path-help` after mounting the secrets
engine.

> [!NOTE]
> **Version information** This documentation assumes you are using at least
> Consul version 1.4 (released 2018). Older version are still supported by the
> plugin, but support is deprecated.

## Quick start

The first step to using the OpenBao secrets engine is to enable it.

```shell-session
$ bao secrets enable consul
Successfully mounted 'consul' at 'consul'!
```

For a quick start, you can use the SecretID token provided by the Consul ACL
bootstrap process, although this is discouraged for production deployments.

```shell-session
$ consul acl bootstrap
AccessorID:       eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee
SecretID:         ffffffff-ffff-ffff-ffff-ffffffffffff
Description:      Bootstrap Token (Global Management)
Local:            false
Create Time:      2025-01-01 10:10:10.124356789 +0000 UTC
Policies:
   00000000-0000-0000-0000-000000000001 - global-management
```

The suggested pattern is to generate a token specifically for OpenBao, following the
[Consul ACL Documentation](https://developer.hashicorp.com/consul/docs/secure/acl)

Next, we must configure OpenBao to know how to contact Consul.
This is done by writing the access information:

```shell-session
$ bao write consul/config/access \
    address=consul.example.com \
    scheme=https \
    token=ffffffff-ffff-ffff-ffff-ffffffffffff
```

In this case, we've configured OpenBao to connect to Consul at
`https://consul.example.com`. We've also provided an ACL token to use with the
`token` parameter. OpenBao must have a token with the permissions to create and
revoke ACL tokens.

The next step is to configure a role. A role is a logical name that maps
to a set of policy names used to generate those credentials. For example, let's create
a "worker" role that maps to a "readonly" policy:

```shell-session
$ bao write consul/roles/worker policies=readonly
```

The secrets engine expects either a single or a comma separated list of policy names.

To generate a new Consul ACL token, we simply read from that role:

```shell-session
$ bao read consul/creds/worker
Key                 Value
---                 -----
lease_id            consul/creds/worker/0aB0aB0aB0aB0aB0aB0aB0aB
lease_duration      1h
lease_renewable     true
accessor            aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
consul_namespace    n/a
local               false
partition           n/a
token               bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb
```

Here we can see that OpenBao has generated a new Consul ACL token for us.
We can test this token out, by reading it in Consul (by it's accessor):

```shell-session
$ consul acl token read -accessor-id aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
AccessorID:       aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
SecretID:         bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb
Description:      Vault worker root 1735726200124356789
Local:            false
Create Time:      2025-01-01 10:10:10.124356789 +0000 UTC
Policies:
   cccccccc-cccc-cccc-cccc-cccccccccccc - readonly
```

## API

The Consul secrets engine has a full HTTP API. Please see the [API
docs](api/readme.md) for more details.
