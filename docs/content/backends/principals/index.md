---
title: Principals
---

## Local

### Example Usage

```
principalsType: local
principalsOpts:
  users:
    foouser: fooprincpal,anotherprincipal,thirdprincipal
    baruser: anotherprincipal
```

### Options

  * **users** - Map of users and associated principals (required)

## LDAP

### Example Usage

```
principalsType: ldap
principalsOpts:
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
