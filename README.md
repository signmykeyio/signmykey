![Signmykey logo](docs/content/images/logo-full.png)

----

![Build Status](https://github.com/signmykeyio/signmykey/actions/workflows/master.yml/badge.svg)

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
curl -Lo signmykey https://github.com/signmykeyio/signmykey/releases/download/v0.5.1/signmykey_linux_amd64
```
* Install it in your PATH:
```
chmod +x signmykey && sudo mv signmykey /usr/bin/
```

### APT

* Ensure you have curl and gpg
```
sudo apt update && sudo apt install ca-certificates curl gnupg
```
* Add Signmykey GPG to your APT truststore
```
curl https://gpg.signmykey.io/signmykey.pub | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/signmykey.gpg
```
* Add Signmykey repository
```
echo 'deb [signed-by=/etc/apt/trusted.gpg.d/signmykey.gpg] https://apt.signmykey.io stable main' | sudo tee /etc/apt/sources.list.d/signmykey.list
```
* Install Signmykey package
```
sudo apt update && sudo apt install signmykey
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
sudo yum install signmykey
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
