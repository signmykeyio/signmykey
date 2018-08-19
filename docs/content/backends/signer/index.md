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
