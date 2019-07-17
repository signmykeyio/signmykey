---
title: SSH Server
---

## Configuration

In order to use SSH principals, you must configure your SSH servers to use them.

You can find [here](/getting-started/vault/#export-ca-public-key) how to generate the file */etc/ssh/ca.pem*.

Modify the file */etc/ssh/sshd_config* with the following parameters

```
...
# Allow root to connect
PermitRootLogin yes

# Copy Vault SSH CA
TrustedUserCAKeys /etc/ssh/ca.pem

# Permit user principals
AuthorizedPrincipalsFile /etc/ssh/authorized_principals/%u

# Deny non signed key files
AuthorizedKeysFile /dev/null

# Deny password authentication
PasswordAuthentication no
...
``` 

## Principals


Create the */etc/ssh/authorized_principals* directory
```sh
mkdir /etc/ssh/authorized_principals/
```

Also create the file */etc/ssh/authorized_principals/root* for the **root** user
```
hackers
superheros
```
It means that users with **hackers** and **superheros** principals can login as **root** to the server with ssh.

## Restart

{{< warning title="Warning" >}}
Be sure to be able to connect via a console to your server.
{{< /warning >}}

```sh
systemctl restart sshd.service
```