---
title: Signmykey
weight: 0
---

## Installation

On Ubuntu 16.04, add this to the file */etc/apt/sources.list.d/signmykey.list*
```
deb https://apt.signmykey.io/signmykey/ ubuntu main
``` 

Then 

```sh
useradd --no-create-home -s /bin/false signmykey
apt update && apt install signmykey
wget https://gitlab.com/signmykey/signmykey/raw/master/signmykey.service -O /etc/systemd/system/signmykey.service
systemctl enable signmykey.service
``` 

## Server configuration

```
mkdir /etc/signmykey
```

File */etc/signmykey/server.yml*

In this file, you put the Vault AppRole credentials.

```
authenticatorType: ldap
authenticatorOpts:
  ldapAddr: localhost
  ldapPort: 3893
  ldapTls: False
  ldapTlsVerify: False
  ldapBindUser: "cn=serviceuser,ou=svcaccts,dc=glauth,dc=com"
  ldapBindPassword: "mysecret" 
  ldapBase: "dc=glauth,dc=com"
  ldapSearch: "(cn=%s)"

principalsType: ldap
principalsOpts:
  ldapAddr: localhost
  ldapPort: 3893
  ldapTls: False
  ldapTlsVerify: False
  ldapBindUser: "cn=serviceuser,ou=svcaccts,dc=glauth,dc=com"
  ldapBindPassword: "mysecret"
  ldapBase: "ou=groups,dc=glauth,dc=com"
  ldapSearch: "(cn=%s)"

signerType: vault
signerOpts:
  vaultAddr: "localhost"
  vaultPort: 8200
  vaultTls: true
  vaultPath: "ssh"
  vaultRole: "sign-user-role"
  vaultRoleid: "140f639f-3c86-4bce-6019-8a9cfd4e47e8"
  vaultSecretid: "959401d7-c66b-8988-0ab3-e56f9ba6a708"
  vaultSignTtl: "24h"
```

Secure the config file

```sh
chown signmykey: /etc/signmykey/config.yml
chmod 600 /etc/signmykey/config.yml
```

### Server start

```sh
systemctl start signmykey.service
systemctl status signmykey.service
``` 

## Client configuration

### Global configuration

```
mkdir /etc/signmykey
```

Content of the */etc/signmykey/client.yml* file:

```
https://signmykeyserver/
```

### User configuration

Content of the *~/.signmykey.yml* file:

```
https://signmykeyserver/
```

## Client usage

### Sign your key

```sh
signmykey
```

### Verify your key principals

```sh
ssh-keygen -Lf ~/.ssh/id_rsa-cert.pub
```

Look at Principals in the output:
```
/home/myuser/.ssh/id_rsa-cert.pub:
        Type: ssh-rsa-cert-v01@openssh.com user certificate
        Public key: RSA-CERT ce:75:5e:4d:3a:db:29:f4:69:3f:98:39:80:48:a3:0f
        Signing CA: RSA 04:cc:f7:15:b6:3a:ab:9a:9a:cf:e8:e4:82:5d:a9:0e
        Key ID: "johndoe"
        Serial: 15992710984477402823
        Valid: from 2018-07-30T15:26:46 to 2018-07-31T15:27:16
        Principals: 
                superheros
                vpn
        Critical Options: (none)
        Extensions: 
                permit-pty
```