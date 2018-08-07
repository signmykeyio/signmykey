---
title: Signmykey
---

## Installation

On Ubuntu 16.04, add this to the file */etc/apt/sources.list.d/signmykey.list*
```
echo "deb [trusted=yes] https://apt.signmykey.io/signmykey/ xenial main" > /etc/apt/sources.list.d/signmykey.list
``` 

Then 

```sh
useradd --no-create-home -s /bin/false signmykey
apt update && apt install signmykey
wget https://gitlab.com/signmykey/signmykey/raw/master/signmykey.service -O /etc/systemd/system/signmykey.service
systemctl enable signmykey.service
``` 

```
mkdir -m 700 /etc/signmykey
```

## Server certificate

Generate a certificate for signmykey server using Vault PKI

**Note**: you can use another certificate provider

```sh
vault write pki/issue/allow-all-domains common_name="signmykeyserver" alt_names="localhost" ip_sans="127.0.0.1"
```

Copy the output from the previous command
```sh
vi /etc/signmykey/server.key # certificate key
vi /etc/signmykey/server.pem # private_key key
chmod 400 /etc/signmykey/server.key
```

## Server configuration

File */etc/signmykey/server.yml*

In this file, you put the Vault AppRole credentials.

```
authenticatorType: ldap
authenticatorOpts:
  ldapAddr: localhost
  ldapPort: 3893
  ldapTLS: False
  ldapTLSVerify: False
  ldapBindUser: "cn=serviceuser,ou=svcaccts,dc=glauth,dc=com"
  ldapBindPassword: "mysecret"
  ldapBase: "dc=glauth,dc=com"
  ldapSearch: "(cn=%s)"

principalsType: ldap
principalsOpts:
  ldapAddr: localhost
  ldapPort: 3893
  ldapTLS: False
  ldapTLSVerify: False
  ldapBindUser: "cn=serviceuser,ou=svcaccts,dc=glauth,dc=com"
  ldapBindPassword: "mysecret"
  ldapBase: "ou=groups,dc=glauth,dc=com"
  ldapSearch: "(cn=%s)"

signerType: vault
signerOpts:
  vaultAddr: "localhost"
  vaultPort: 8200
  vaultTLS: true
  vaultPath: "ssh"
  vaultRole: "sign-user-role"
  vaultRoleID: "11940c2d-4639-9358-d750-cdb7cf409ff4"
  vaultSecretID: "8b4c901f-1f84-5049-17ee-92de12b6b1e5"
  vaultSignTTL: "24h"

address: "0.0.0.0:443"
tlsDisable: false
tlsCert: "/etc/signmykey/server.pem" 
tlsKey: "/etc/signmykey/server.key"
```

Secure the config file

```sh
chmod 600 /etc/signmykey/server.yml
chown -R signmykey: /etc/signmykey
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
addr: "https://signmykeyserver/"
```

### User configuration

Content of the *~/.signmykey.yml* file:

```
addr: "https://signmykeyserver/"
```

## Client usage

### Sign your key

```sh
signmykey -u johndoe
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