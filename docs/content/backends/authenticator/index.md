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
    otpuser: $2a$10$/6T2iN8I7UTUTDuezVH41eDlSIeNr32wi9PtDfNF3Zxes3RO0LK/a,VHOUR7WH7N6ZXI5VEKZFZ4ESB4ZEYPGNUDAT6LKGNHLWUXMTEYKA====
```

### Options

  * **users** - Map of users and bcrypt hashed passwords (you can hash passwords via "signmykey hash" command) (required)

Optionally, users may wish to utilize OTP in which case the "signmykey hash" command generates a longer string which contains the 
hashed password and the OTP seed encrypted with the user's password. 

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

## OIDC ROPC

### Example Usage

```
authenticatorType: oidcropc
authenticatorOpts:
  oidcTokenEndpoint: "https://idp.my.corp/auth/realms/mycorp/protocol/openid-connect/token"
  oidcClientID: "signmykey"
  oidcClientSecret: "93fac2d9-bd8f-453a-9ece-e2c430f0ee04"
```

### Options

  * **oidcTokenEndpoint** - OpenID Connect token Endpoint (required)
  * **oidcClientID** - OpenID Connect Client ID (required)
  * **oidcClientSecret** - OpenID Connect Client Secret (required)
