---
title: Signmykey
type: index
---

## signmykey

### What is it ?

Signmykey is a server and a client to sign ssh keys with Hashicorp Vault.

It helps you to delegate the signing part and the principals in your oganization with the same workflow.

As an example, signmykey allows you to let your LDAP users to sign theirs ssh keys with principals using the **memberOf** attribute.

### What do you need ?

You need :

- an LDAP server 
- a Vault instance
- a signmykey host


### How does it work ?

![Signmykey workflow](images/signmykey.png)

#### Signmykey workflow (blue)

1. The user enters `signmykey` command then inserts its LDAP credentials, Signmykey client push to Signmykey server the credentials and SSH public key
2. The LDAP server verifies the credentials and gives user groups to Signmykey server
3. Signmykey server asks Vault to sign the public key with groups as principals
4. Vault gives back the signed key to Signmykey server
5. Signmykey server gives back the signed key to the user

#### SSH workflow (green)

1. The user enters `ssh root@server1` from his terminal
2. The SSH server verifies that the certificate 
    - is signed with the CA
    - is valid
    - has the correct principals in /etc/ssh/authorized_principals/root