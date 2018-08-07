---
title: LDAP Server
---

This HOWTO describe how to configure an LDAP server using [glauth](https://github.com/glauth/glauth).

## Installation

### Download
```sh
wget https://github.com/glauth/glauth/releases/download/v1.1.0/glauth64 -O /usr/local/bin/glauth
chmod +x /usr/local/bin/glauth
```

### User and directories

```sh
useradd --no-create-home -s /bin/false glauth
mkdir -m 700 /etc/glauth
wget https://raw.githubusercontent.com/glauth/glauth/master/sample-simple.cfg -O /etc/glauth/glauth.cfg
chown -R glauth: /etc/glauth/
```
### Configuration

The default configuration downloaded during the previous step contains some users.

### Launch

```sh
runuser -s /bin/bash -l glauth -c 'glauth -c /etc/glauth/glauth.cfg'
```

## Usage

List LDAP entries

```sh
ldapsearch -LLL -H ldap://localhost:3893 -D cn=serviceuser,ou=svcaccts,dc=glauth,dc=com -w mysecret -x -bdc=glauth,dc=com
```