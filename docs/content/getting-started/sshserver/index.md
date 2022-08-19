
# SSH Server


## Configuration

In order to use SSH principals, you must configure your SSH servers to use them.

You can find [here](/getting-started/vault/#export-ca-public-key) how to generate the file */etc/ssh/ca.pem*.


### Linux server

{{< warning title="Warning" >}} Open SSH Server must be installed on the server. {{< /warning >}}
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

#### Principals

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

#### Restart

{{< warning title="Warning" >}}
Be sure to be able to connect via a console to your server.
{{< /warning >}}

```sh
systemctl restart sshd.service
```


### Windows

{{< warning title="Warning" >}} Open SSH Server must be installed on the server. On Windows 2019 Server the service is present by default, just enable it. {{< /warning >}}
Modify the file *C:\ProgramData\ssh\sshd_config* with the following parameters

```
...
# Allow log into **
SyslogFacility LOCAL0

# Copy Vault SSH CA
TrustedUserCAKeys __PROGRAMDATA__/ssh/trusted-vault-ca-keys.pub

# Permit user principals
AuthorizedPrincipalsFile __PROGRAMDATA__/ssh/authorized_principals

# Deny non signed key files
AuthorizedKeysFile  none

# Deny password authentication
PasswordAuthentication no

# Nor mandatory but usefull when user has Administrators' right
Match Group administrators
#       AuthorizedKeysFile __PROGRAMDATA__/ssh/administrators_authorized_keys
       AuthorizedPrincipalsFile __PROGRAMDATA__/ssh/authorized_principals
...
```

#### Principals

Create the *C:\ProgramData\ssh\authorized_princpals* file
Be sure to only grant RO on this file for SYSTEM user.

In this file indicate all the principals that are owned by users you want to be able to connect to the server.
```
hackers, superheros
```
It means that users with **hackers** and **superheros** principals can login to the server with ssh.

#### Restart

{{< warning title="Warning" >}}
Be sure to be able to connect via a console/rdp to your server.
{{< /warning >}}

Reload the service in the in the *Services* application
or
In cmd.exe
```
net stop gsw_sshd && net start gsw_sshd
```

### Additionnal infos

#### Deprecated algorithm

In recent OS (Ubuntu 22.04), not yet on Redhat family (9) you can no longer ssh to servers with that OS because the RSA algorithm is deprecated.
Therefor you have two choices :

##### Update signmykey and regenerate key

Check signmykey page
```
[signmykey](../../signmykey/index.md)
```

##### Authorize deprecated algorithm

On the server add this to the ssh configuration file :
```
PubkeyAcceptedAlgorithms +ssh-rsa-cert-v01@openssh.com
```

Restart ssh service :
```
systemctl restart sshd
```
