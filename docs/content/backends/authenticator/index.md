---
title: Authenticator
---

## LDAP

### Example Usage

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
```

### Options

  * **ldapAddr** - Address of LDAP server (required)
  * **ldapPort** - Port of LDAP server
  * **ldapTLS** - Enable/disable SSL/TLS connection
  * **ldapTLSVerify** - Enable/disable verification of SSL/TLS certificate
  * **ldapBindUser** - LDAP bind user
  * **ldapBindPassword** - LDAP bind password
  * **ldapBase** - LDAP search base
  * **ldapSearch** - LDAP search string to find user

## Vault

### Example Usage

```
```

### Options

  * **vaultAddr** - Address of Vault server
  * **vaultPort** - Port of Vault server
  * **vaultTLS** - Enable/disable SSL/TLS connection to Vault server
  * **vaultPath** - Authentication path on Vault server of userpass auth backend
