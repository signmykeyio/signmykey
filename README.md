![Signmykey logo](docs/content/images/logo-full.png)

----

[![Build Status](https://travis-ci.org/signmykeyio/signmykey.svg?branch=master)](https://travis-ci.org/signmykeyio/signmykey) [![Go Report Card](https://goreportcard.com/badge/github.com/signmykeyio/signmykey)](https://goreportcard.com/report/github.com/signmykeyio/signmykey) [![Maintainability](https://api.codeclimate.com/v1/badges/bc6e89d9e4d60b2d688f/maintainability)](https://codeclimate.com/github/signmykeyio/signmykey/maintainability)

----

Signmykey is an automated SSH Certificate Authority. It allows you to securly and centraly manage SSH accesses to your infrastructure.

Three types of backends are supported by Signmykey:

* **Authorization**: users can be authentified through different systems like LDAP or Local map.
* **Principals**: list of principals applied to SSH certificates can be created dynamically from LDAP groups or set staticaly in local config.
* **Signer**: cryptographic signing operations of SSH certificates can be done directly by Signmykey or via Hashicorp Vault.

## Install

* Download **signmykey** binary (on 64bits linux):
```
curl -Lo signmykey https://github.com/signmykeyio/signmykey/releases/download/0.2.1/signmykey_linux_amd64
```
* Add execute permission:
```
chmod +x signmykey
```
* Install it in your PATH:
```
sudo mv signmykey /usr/bin/signmykey
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
