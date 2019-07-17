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
  ldapUserBase: "dc=glauth,dc=com"
  ldapUserSearch: "(cn=%s)"
  ldapGroupBase: "dc=glauth,dc=com"
  ldapGroupSearch: "(&(objectClass=group)((member=%s)))"
```

### Options

  * **ldapAddr** - Address of LDAP server (required)
  * **ldapPort** - Port of LDAP server
  * **ldapTLS** - Enable/disable SSL/TLS connection
  * **ldapTLSVerify** - Enable/disable verification of SSL/TLS certificate
  * **ldapBindUser** - LDAP bind user
  * **ldapBindPassword** - LDAP bind password
  * **ldapUserBase** - LDAP user search base
  * **ldapUserSearch** - LDAP search string to find user
  * **ldapGroupBase** - LDAP groups search base
  * **ldapGroupSearch** - LDAP search string to find groups
  * **ldapGroupPrefix** - Filter LDAP groups by prefix
  * **transformCase** - Change case of returned principals (default: none) (must be "none", "lower" or "upper")

## OIDC ROPC

### Example Usage

```
principalsType: oidcropc
principalsOpts:
  oidcUserinfoEndpoint: "https://idp.my.corp/auth/realms/mycorp/protocol/openid-connect/userinfo"
  oidcUserGroupsEntry: "oidc-groups"
  transformCase: upper
```

### Options

  * **oidcUserinfoEndpoint** - OpenID Connect userinfo Endpoint (required)
  * **oidcUserGroupsEntry** - OpenID Connect group entry name returned by userinfo endpoint (required)
  * **transformCase** - Change case of returned principals (default: none) (must be "none", "lower" or "upper")
