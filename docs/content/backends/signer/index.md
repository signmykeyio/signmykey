---
title: Signer
---

## Local

### Example Usage

```
signerType: local
signerOpts:
  caCert: /etc/signmykey/ca.pub
  caKey: /etc/signmykey/ca
  ttl: 300
  extensions:
    permit-pty: ""
```

### Options

  * **caCert** - Path to CA public key (required)
  * **caKey** - Path to CA private key (required)
  * **ttl** - TTL in seconds for signed certificates (required)
  * **criticalOptions** - Map of critical options for signed certificates (optional) (default: empty)
  * **extensions** - Map of extensions for signed certificates (optional) (default: permit-X11-forwarding, permit-agent-forwarding, permit-port-forwarding, permit-pty, permit-user-rc)

## Vault (Hashicorp)

### Example Usage

```
signerType: vault
signerOpts:
  vaultAddr: 127.0.0.1
  vaultPort: 8200
  vaultTLS: true
  vaultRoleID: db02de05-fa39-4855-059b-67221c5c2f63
  vaultSecretID: 6a174c20-f6de-a53c-74d2-6018fcceff64
  vaultPath: ssh
  vaultRole: ssh-client
  vaultSignTTL: 600
```

### Options

  * **vaultAddr** - Address of Vault server
  * **vaultPort** - Port of Vault server
  * **vaultTLS** - Enable/disable SSL/TLS connection to Vault server
  * **vaultRoleID** - Approle Role ID to connect to Vault
  * **vaultSecretID** - Approle Secret ID to connect to Vault
  * **vaultPath** - Path to SSH Signed certificates secret backend on Vault server
  * **vaultRole** - Role of SSH secret backend to use for ssh key signing
  * **vaultSignTTL** - TTL to apply to signed keys
