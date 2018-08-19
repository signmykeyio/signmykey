## 0.1.1 (August 19th, 2018)

BUG FIXES:
  * server:
    * fix dev mode #30

## 0.1.0 (August 19th, 2018)

DEPRECATIONS/CHANGES:
  * server:
    * all config keys are now camelCase instead of snake-case
    * default listening port is now 9600 instead of 8080 #14
    * local principals usermap is now under "users" subkey #27

FEATURES:
  * add documentation at https://signmykey.io/
  * server:
    * Local Signer: new backend to able to sign SSH keys without Hashicorp Vault #11
    * Local Authenticator: new backend to able to Authenticate users without LDAP #12
    * Dev mode: add a dev/demo in-memory mode to start signmykey server without config #15
    * RPM: add rpm package creation in Makefile #26

IMPROVEMENTS:
  * server:
    * enable HTTPS support to expose signmykey API
    * add timeout on ldap and vault connections
