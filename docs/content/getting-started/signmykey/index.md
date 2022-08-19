---
title: Signmykey
---

## Installation

On Ubuntu 20.04+, add signmykey repository and key:
```
echo "deb https://apt.signmykey.io/ stable main" > /etc/apt/sources.list.d/signmykey.list
curl https://gpg.signmykey.io/signmykey.pub | apt-key add -
```

Then

```sh
useradd --no-create-home -s /bin/false signmykey
apt update && apt install signmykey
wget https://raw.githubusercontent.com/signmykeyio/signmykey/master/signmykey.service -O /etc/systemd/system/signmykey.service
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

### LDAP configuration alternative

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
  vaultSignTTL: "12h"

address: "0.0.0.0:443"
tlsDisable: false
tlsCert: "/etc/signmykey/server.pem" 
tlsKey: "/etc/signmykey/server.key"
```

### OIDC ROPC configuration alternative
```
authenticatorType: oidcropc
authenticatorOpts:
  oidcTokenEndpoint: "https://idp.my.corp/auth/realms/mycorp/protocol/openid-connect/token"
  oidcClientID: "signmykey"
  oidcClientSecret: "93fac2d9-bd8f-453a-9ece-e2c430f0ee04"

principalsType: oidcropc
principalsOpts:
  oidcUserinfoEndpoint: "https://idp.my.corp/auth/realms/mycorp/protocol/openid-connect/userinfo"
  oidcUserGroupsEntry: "oidc-groups"

signerType: vault
signerOpts:
  vaultAddr: "localhost"
  vaultPort: 8200
  vaultTLS: true
  vaultPath: "ssh"
  vaultRole: "sign-user-role"
  vaultRoleID: "11940c2d-4639-9358-d750-cdb7cf409ff4"
  vaultSecretID: "8b4c901f-1f84-5049-17ee-92de12b6b1e5"
  vaultSignTTL: "12h"

address: "0.0.0.0:443"
tlsDisable: false
tlsCert: "/etc/signmykey/server.pem" 
tlsKey: "/etc/signmykey/server.key"
```

### Secure the config file

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

## Additionnal infos

### Deprecated algorithm

In recent OS (Ubuntu 22.04), not yet on Redhat family (9) you can no longer ssh to servers with that OS because the RSA algorithm is deprecated.
Therefor you have two choices :

#### Update signmykey and regenerate key
Download the new version of signmykey that support every ssh keygen algorithm (>= 0.7.0)
```
smkVersion='0.7.0'
wget https://github.com/signmykeyio/signmykey/releases/download/v$smkVersion/signmykey_linux_amd64 -O /tmp/signmykey
chmod +x /tmp/signmykey
smkPath=`which signmykey | sed 's|\(.*\)/.*|\1|'`
mv /tmp/signmykey $smkPath
signmykey version
```

Generate a new sshkey with a supported algorithm
```
ssh-keygen -t ed25519
```

Sign you key
```
signmykey -u johndoe
```

#### Authorize deprecated algorithm

Check sshserver page
```
[sshserver](/signmykey/docs/content/getting-started/sshserver/index.md)
```
