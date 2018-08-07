---
title: Vault PKI
---

This HOWTO describe how to configure a Vault server with a PKI to sign its own certificate.

## Installation

### Download
```sh
apt install unzip
VAULTVERSION=0.10.4
cd /tmp
wget https://releases.hashicorp.com/vault/${VAULTVERSION}/vault_${VAULTVERSION}_linux_amd64.zip
unzip vault_${VAULTVERSION}_linux_amd64.zip && rm vault_${VAULTVERSION}_linux_amd64.zip
sudo mv vault /usr/local/bin/vault
```

### User and directories

```sh
useradd --no-create-home -s /bin/false vault
mkdir -m 700 /etc/vault /var/lib/vault
chown vault: /etc/vault /var/lib/vault
```

### Configuration

Config file */etc/vault/vault.conf*
```
backend "file" {
  path = "/var/lib/vault"
}
 
listener "tcp" {
  address = "0.0.0.0:8200"
  #tls_cert_file = "/etc/vault/vault.pem"
  #tls_key_file = "/etc/vault/vault.key"
  tls_disable = 1
}

disable_mlock = true
ui = true
```

### Launch
```
runuser -s /bin/bash -l vault -c 'vault server -config /etc/vault/vault.conf'
```

### Init
```sh
export VAULT_ADDR=http://127.0.0.1:8200/
vault operator init -key-shares=1 -key-threshold=1
```

Keys
```
Unseal Key 1: p8zOywmj882TWCfIX0xekQJQcHuGQxm4QwmUyT7qDkA=

Initial Root Token: 5682b22a-02cd-4678-7883-202afed596dc
```

### Unseal

```sh
vault operator unseal
```

### vault root login
```sh
vault login
```

## Enable PKI 

```
vault secrets enable pki
vault secrets tune -max-lease-ttl=87600h pki
vault write pki/root/generate/internal common_name=MyCA ttl=87600h
Key              Value
---              -----
certificate      -----BEGIN CERTIFICATE-----
MIIDIDCCAgigAwIBAgIUSzV637STZohgSiOMrvDD9Ke1/5owDQYJKoZIhvcNAQEL
BQAwDzENMAsGA1UEAxMETXlDQTAeFw0xODA4MDYxMjE3NDhaFw0yODA4MDMxMjE4
MThaMA8xDTALBgNVBAMTBE15Q0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDBlBgJna7MaKUV5+q3L5edm0eKsZomR6NCUUNMfd02IuahtSWrHvFp4vSV
UL9x0jX5fxdaez1KPntm1Nd6sILW5jCr1weZ4oPuVn6E33qyVGFDg/8w8WL46VRc
PeuoD8/zarze71gvUKRnQ/QHu6IJ4FzRNexX78BiFcvocrH/ci3LjEXZWOHSsXs5
yRMEwQFOVTQnek2Ui6r1u9fe+YEzGMvW2NcSp6lF6/KMYCULqfzsS2JuEUxQOmh4
z9yyHVD0Eq1OsRsdo4yTS+5z3LlkzVhoBQWBVgRcJJSNTfZKzfnOipFSTUXs3Yvp
DrHnBH3zBb0NBqmL+AZDL393QlSrAgMBAAGjdDByMA4GA1UdDwEB/wQEAwIBBjAP
BgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTiJ/SUuFBhMUKmJrDYlihZpIjTbzAf
BgNVHSMEGDAWgBTiJ/SUuFBhMUKmJrDYlihZpIjTbzAPBgNVHREECDAGggRNeUNB
MA0GCSqGSIb3DQEBCwUAA4IBAQBlst64whBAry9mOz/xso0SogeqKsRUCbjUYegH
6+XigZJcyBkKioJwVxJ9m625R3JGrS8xC3U0R7IHX3f0+3nR8Q6+DkNKK+mXCPoa
+8Rw5kEJDW94v86stUPsBblU+Y97Zva4mtog73k+aoWX9Ok39j3VE4jCE1126bGw
nvteK/ghIQDQ4G9JQvFJMrT+iWMP5LiduYdy4ANUifomGRGS8YHSN0ae+S7ICPm3
GCUb1fAIt02FdUk7U+bV/6a8XPx1e8lqKpopM+ylflHS02bcTXlzwYOV1x9aAPp8
LgBMtfqPzmS4UxfWhPdhaF+48rA2b5zGXEy5Lj7Lr+JcAMIC
-----END CERTIFICATE-----
expiration       1848917898
issuing_ca       -----BEGIN CERTIFICATE-----
MIIDIDCCAgigAwIBAgIUSzV637STZohgSiOMrvDD9Ke1/5owDQYJKoZIhvcNAQEL
BQAwDzENMAsGA1UEAxMETXlDQTAeFw0xODA4MDYxMjE3NDhaFw0yODA4MDMxMjE4
MThaMA8xDTALBgNVBAMTBE15Q0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDBlBgJna7MaKUV5+q3L5edm0eKsZomR6NCUUNMfd02IuahtSWrHvFp4vSV
UL9x0jX5fxdaez1KPntm1Nd6sILW5jCr1weZ4oPuVn6E33qyVGFDg/8w8WL46VRc
PeuoD8/zarze71gvUKRnQ/QHu6IJ4FzRNexX78BiFcvocrH/ci3LjEXZWOHSsXs5
yRMEwQFOVTQnek2Ui6r1u9fe+YEzGMvW2NcSp6lF6/KMYCULqfzsS2JuEUxQOmh4
z9yyHVD0Eq1OsRsdo4yTS+5z3LlkzVhoBQWBVgRcJJSNTfZKzfnOipFSTUXs3Yvp
DrHnBH3zBb0NBqmL+AZDL393QlSrAgMBAAGjdDByMA4GA1UdDwEB/wQEAwIBBjAP
BgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTiJ/SUuFBhMUKmJrDYlihZpIjTbzAf
BgNVHSMEGDAWgBTiJ/SUuFBhMUKmJrDYlihZpIjTbzAPBgNVHREECDAGggRNeUNB
MA0GCSqGSIb3DQEBCwUAA4IBAQBlst64whBAry9mOz/xso0SogeqKsRUCbjUYegH
6+XigZJcyBkKioJwVxJ9m625R3JGrS8xC3U0R7IHX3f0+3nR8Q6+DkNKK+mXCPoa
+8Rw5kEJDW94v86stUPsBblU+Y97Zva4mtog73k+aoWX9Ok39j3VE4jCE1126bGw
nvteK/ghIQDQ4G9JQvFJMrT+iWMP5LiduYdy4ANUifomGRGS8YHSN0ae+S7ICPm3
GCUb1fAIt02FdUk7U+bV/6a8XPx1e8lqKpopM+ylflHS02bcTXlzwYOV1x9aAPp8
LgBMtfqPzmS4UxfWhPdhaF+48rA2b5zGXEy5Lj7Lr+JcAMIC
-----END CERTIFICATE-----
serial_number    4b:35:7a:df:b4:93:66:88:60:4a:23:8c:ae:f0:c3:f4:a7:b5:ff:9a
```

### Add CA public cert to server

```sh
wget http://127.0.0.1:8200/v1/pki/ca/pem -O /usr/local/share/ca-certificates/MyCA.crt
update-ca-certificates
```

### Add CA public cert to firefox
https://www.cyberciti.biz/faq/firefox-adding-trusted-ca/


### CA role

```sh
vault write pki/roles/allow-all-domains allow_any_name=true max_ttl=8760h
```

### Generate vault own certificate

```sh
vault write pki/issue/allow-all-domains common_name="vaultserver" alt_names="localhost" ip_sans="127.0.0.1"
Key                 Value
---                 -----
certificate         -----BEGIN CERTIFICATE-----
MIIDTzCCAjegAwIBAgIUc6DXr149vhgqzYB6W6FhK65vbgAwDQYJKoZIhvcNAQEL
BQAwDzENMAsGA1UEAxMETXlDQTAeFw0xODA4MDYxMjMwMDlaFw0xODA5MDcxMjMw
MzlaMBYxFDASBgNVBAMTC3ZhdWx0c2VydmVyMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAsPwTd/opXlgW8kCiLDvnUxcUWcMxTEqvALuAtaYhbouK7OBI
pt8cT7lZ91VAH6UAKFh7eQqRXjLpHmEz6qfMNuIN+RShUvanTRFht/urY8Z/Be9/
wpcL/vJuCh8Y65gqI9I9dj5NmayGJ71v7aoKoofPnR5j2J2qFHANXn4VX+51St0H
FbhU2xIftCOB3nF4wtvZhOeXde4nAphOmjdfEvQNFjjdXZkqfvlR3N/ta7v+3At+
VXzp5LWLbLvpe80v2otqR0G/Ja29rVYLLL8MM0FCyjtOFMj4ezWW0C8FQ/+YZ02g
9FBNL3aSEKqsE0nPsVc2ijiJB0moqVQBJoj3rQIDAQABo4GbMIGYMA4GA1UdDwEB
/wQEAwIDqDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0OBBYE
FCKx0/NQNxxR16oN4piraGaL82+iMB8GA1UdIwQYMBaAFOIn9JS4UGExQqYmsNiW
KFmkiNNvMCcGA1UdEQQgMB6CC3ZhdWx0c2VydmVygglsb2NhbGhvc3SHBH8AAAEw
DQYJKoZIhvcNAQELBQADggEBAEs84gw0k9DODkhSAae4a73LHCWLiyjwN41Y+RCA
/9+OAJrLbfX/IHP/J+C9F7KTHcV+dpisW5yAcdtDV1VwyF9I8np5Jnp1h/ry1jVl
3YvO22cVx57bsU+DUcNBCWHWDCESGSlJThHIvaTCqPwobFEnDOGv9N+ITWfcc+hJ
BIkevpb9W7iCxOV8KMBBsJQrhClVkwHsik7sODESWG+fWLIuZh/fVTBOegKmVazb
+Q3wYoHja7+osoqcUssLPro1/SYCunDt1UI/qKU43/JgDTFXTBzSqCiqG6QU7XHg
Max9HpPB/0uDtQjsr8VVOw1/ST4dDLdOBT7I2/0LcyhzuCo=
-----END CERTIFICATE-----
issuing_ca          -----BEGIN CERTIFICATE-----
MIIDIDCCAgigAwIBAgIUSzV637STZohgSiOMrvDD9Ke1/5owDQYJKoZIhvcNAQEL
BQAwDzENMAsGA1UEAxMETXlDQTAeFw0xODA4MDYxMjE3NDhaFw0yODA4MDMxMjE4
MThaMA8xDTALBgNVBAMTBE15Q0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDBlBgJna7MaKUV5+q3L5edm0eKsZomR6NCUUNMfd02IuahtSWrHvFp4vSV
UL9x0jX5fxdaez1KPntm1Nd6sILW5jCr1weZ4oPuVn6E33qyVGFDg/8w8WL46VRc
PeuoD8/zarze71gvUKRnQ/QHu6IJ4FzRNexX78BiFcvocrH/ci3LjEXZWOHSsXs5
yRMEwQFOVTQnek2Ui6r1u9fe+YEzGMvW2NcSp6lF6/KMYCULqfzsS2JuEUxQOmh4
z9yyHVD0Eq1OsRsdo4yTS+5z3LlkzVhoBQWBVgRcJJSNTfZKzfnOipFSTUXs3Yvp
DrHnBH3zBb0NBqmL+AZDL393QlSrAgMBAAGjdDByMA4GA1UdDwEB/wQEAwIBBjAP
BgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTiJ/SUuFBhMUKmJrDYlihZpIjTbzAf
BgNVHSMEGDAWgBTiJ/SUuFBhMUKmJrDYlihZpIjTbzAPBgNVHREECDAGggRNeUNB
MA0GCSqGSIb3DQEBCwUAA4IBAQBlst64whBAry9mOz/xso0SogeqKsRUCbjUYegH
6+XigZJcyBkKioJwVxJ9m625R3JGrS8xC3U0R7IHX3f0+3nR8Q6+DkNKK+mXCPoa
+8Rw5kEJDW94v86stUPsBblU+Y97Zva4mtog73k+aoWX9Ok39j3VE4jCE1126bGw
nvteK/ghIQDQ4G9JQvFJMrT+iWMP5LiduYdy4ANUifomGRGS8YHSN0ae+S7ICPm3
GCUb1fAIt02FdUk7U+bV/6a8XPx1e8lqKpopM+ylflHS02bcTXlzwYOV1x9aAPp8
LgBMtfqPzmS4UxfWhPdhaF+48rA2b5zGXEy5Lj7Lr+JcAMIC
-----END CERTIFICATE-----
private_key         -----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAsPwTd/opXlgW8kCiLDvnUxcUWcMxTEqvALuAtaYhbouK7OBI
pt8cT7lZ91VAH6UAKFh7eQqRXjLpHmEz6qfMNuIN+RShUvanTRFht/urY8Z/Be9/
wpcL/vJuCh8Y65gqI9I9dj5NmayGJ71v7aoKoofPnR5j2J2qFHANXn4VX+51St0H
FbhU2xIftCOB3nF4wtvZhOeXde4nAphOmjdfEvQNFjjdXZkqfvlR3N/ta7v+3At+
VXzp5LWLbLvpe80v2otqR0G/Ja29rVYLLL8MM0FCyjtOFMj4ezWW0C8FQ/+YZ02g
9FBNL3aSEKqsE0nPsVc2ijiJB0moqVQBJoj3rQIDAQABAoIBADYUsPZGcQGtNUXN
TkDqBSO0t7k+FgBUCenVYd0f6LNY3JjJaCnln0cVQlJ7sF57EvNBJmm0Ovtn3ygz
V6PqplJW/SIRlcI+MJ0yJIQN2S9h5kqwBoA1m6rJ9aoOGpVTJ/1OLw3Et/2vZEcc
celTvnAvw9clBCma8+/O1ab8LBflycVsakLJMVFVPBrc7kJ3jprqWNi+WbjRecAs
ndlZfXpBk8UiN+ZL5DGP9b/qHsaT+QwmCRhjaOkrh3LF0ah7G7AETdrSZGfXnmhw
MjxXPZkdVp+JYHeeKwFq69xNC/drysV0OooA21gYHXVE4EmB9pnajE2OtaX6p2oY
gnvXH5kCgYEA0oOsrPSeHfk8FjuCMF9h0Iapce/0nmZm7a0jFhJaJ8vIiU1wBu1a
ld8BJn4akoSQLaUZvEzidTYc/EGEoRFZT3GZz3bobwGWdcwNAMuSm3qH0MncWOsV
8yIx97uBWk0HDPaOwmYrI756nfQIvGX1043IV9ThR4vJ6G8Xmc+eD/sCgYEA1znA
TAbTphdDeE1H/9F3FiqLLOrwvh1EkGsvWLg+sQEuJdfK2B6JJw1+1p40JxH8ZUu3
X4T6rCeL8mpF4EDE3rVJRlo6Lr9vLRXXxWtOAqBl/65nZQfDmxJs3Qn62bhpfwZM
YaidTvYQ7I1ND20TWcrTpH3SjIMWp1domef2fncCgYBYnDhY7PaJY1mZeh8IwX1o
yuYUIY70BeKZdOFp7vun+K1GriPTpqEUqLPRQg9pUQdnTzGQA0TnVYnJ3MI5EhZn
zEeT/ldEMoTkvKlUhlwFugPlLLLlcr7ggqpJvtFp8zZejIH27g6Gky0Fw6zRsJFT
JUEJR4A0H3Ezt19VzQCZdQKBgETTpRkq/bgZrGvmWuYGOE0QYd2FbGN/vJNqk4ON
uA6mz/kuHyIp8bZZbHx5rzfnWo2SPxv+zKMNKoXlUl86lzqZQsuKwxx7/7OtTolF
nXbdkIDJZys55mXK6KFvNZc2kBYdD4QThergad0b+s66FPwcDr6FtjVVHoN5Qmwl
cABVAoGAbYHJTrUF8lLz/QyFBVux9M8+w9GGfcXy4FfwPb41QV7cIqpph0s340Pu
8RMz6JMlZrtU5GOUF+sGg5uDpxwatjYThjvZz4dJ4yHmd1JuX0G+jAEBnKf3bYbB
YCP2+bQiqSG4BwzR3gT1iuQZv0OUW6UIYFY3gD8Gg2Cp/q5GHFc=
-----END RSA PRIVATE KEY-----
private_key_type    rsa
serial_number       73:a0:d7:af:5e:3d:be:18:2a:cd:80:7a:5b:a1:61:2b:ae:6f:6e:00
```

### Create cert files from previous step
```
vi /etc/vault/vault.key # private_key entry from previous step
chmod 400 /etc/vault/vault.key
vi /etc/vault/vault.pem # certificate entry from previous step
chown vault: /etc/vault/vault.key /etc/vault/vault.pem
```
### Vault config with TLS

First stop vault server then modify its configuration */etc/vault/vault.conf*

```
backend "file" {
  path = "/var/lib/vault"
}
 
listener "tcp" {
  address = "0.0.0.0:8200"
  tls_cert_file = "/etc/vault/vault.pem"
  tls_key_file = "/etc/vault/vault.key"
  #tls_disable = 1
}

disable_mlock = true
ui = true
```

Then start vault server again
```
runuser -s /bin/bash -l vault -c 'vault server -config /etc/vault/vault.conf'

export VAULT_ADDR=https://localhost:8200/
```

### Then unseal again
```sh
vault status
vault operator unseal
```

