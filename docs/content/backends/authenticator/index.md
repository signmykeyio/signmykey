---
title: Authenticator
---

## Local

### Example Usage

```
authenticatorType: local
authenticatorOpts:
  users:
    foouser: $2a$10$zsvMZ7nEYo4jJJxgb5FpH.izPH37LsuLBXPbuKH4MPF4sihFSG6bW
    baruser: $2a$10$srGqC9g46xaRXbueLk5kDuSuDM6h2EpC.MTRiVaij6s/jcsKQ6LHu
```

### Options

  * **users** - Map of users and bcrypt hashed passwords (you can hash passwords via "signmykey hash" command) (required)

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
