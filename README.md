:warning: **APT GPG Key deprecation** :warning:

The GPG key used to sign APT repository expired. A new one is released at gpg.signmykey.io. To add it to your APT truststore, use this command:
```
curl https://gpg.signmykey.io/signmykey.pub | sudo apt-key add -
```
:warning: **APT GPG Key deprecation** :warning:

----

![Signmykey logo](docs/content/images/logo-full.png)

----

[![Build Status](https://travis-ci.org/signmykeyio/signmykey.svg?branch=master)](https://travis-ci.org/signmykeyio/signmykey) [![Go Report Card](https://goreportcard.com/badge/github.com/signmykeyio/signmykey)](https://goreportcard.com/report/github.com/signmykeyio/signmykey) [![Maintainability](https://api.codeclimate.com/v1/badges/bc6e89d9e4d60b2d688f/maintainability)](https://codeclimate.com/github/signmykeyio/signmykey/maintainability)

----

Signmykey is an automated SSH Certificate Authority. It allows you to securely and centrally manage SSH accesses to your infrastructure.

Three types of backends are supported by Signmykey:

* **Authenticator**: users can be authenticated through different systems like LDAP or Local map.
* **Principals**: list of principals applied to SSH certificates can be created dynamically from LDAP groups or set statically in local config.
* **Signer**: cryptographic signing operations of SSH certificates can be done directly by Signmykey or via Hashicorp Vault.

## Install

### Manual

* Download **signmykey** zip file (ex: on 64bits linux):
```
curl -Lo signmykey https://github.com/signmykeyio/signmykey/releases/download/v0.5.0/signmykey_linux_amd64
```
* Install it in your PATH:
```
chmod +x signmykey && sudo mv signmykey /usr/bin/
```

### APT

* Ensure you have curl and gpg
```
apt update && apt install ca-certificates curl gnupg
```
* Add Signmykey GPG to your APT truststore
```
curl https://gpg.signmykey.io/signmykey.pub | apt-key add -
```
* Add Signmykey repository
```
echo 'deb https://apt.signmykey.io stable main' >> /etc/apt/sources.list.d/signmykey.list
```
* Install Signmykey package
```
apt update && apt install signmykey
```

### YUM

* Add Signmykey repository
```
echo "[signmykey]
name=Signmykey repo
baseurl=https://rpm.signmykey.io/
enabled=1
gpgcheck=0
repo_gpgcheck=1
gpgkey=https://gpg.signmykey.io/signmykey.pub" > /etc/yum.repos.d/signmykey.repo
```
* Install Signmykey package
```
yum install signmykey
```

## Quickstart

* Start server in dev mode (replace *myremoteuser* by the name of the user you want to connect on remote server):
```
signmykey server dev -u myremoteuser
```

* Follow "Server side" instructions displayed by previous command, ex:
```
### Server side                                                                                                                                                                        
                                                                                                                                                                                       
An ephemeral certificate authority is created for this instance and will die with it.                                                                                                  
To deploy this CA on destination servers, you can launch this command:                                                                                                                 
                                                                                                                                                                                       
        $ echo "ssh-rsa fakeCApubKey" > /etc/ssh/ca.pub

You then have to add this line to "/etc/ssh/sshd_config" and restart OpenSSH server:

        TrustedUserCAKeys /etc/ssh/ca.pub
```

* Follow "Client side" instructions, ex:
```
### Client side

A temporary user is created with this parameters:

        user: myremoteuser
        password: fakepassword
        principals: myremoteuser

You can sign your key with this command:

        $ signmykey -a http://127.0.0.1:9600/ -u myremoteuser
```

* Congrats \o/

## Documentation

Documentation is available at https://signmykey.io/
