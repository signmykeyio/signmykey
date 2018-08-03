---
title: Vault
---

Before you can use Signmykey, you must configure Vault.

## Enable ssh

First, enable the ssh secret engine to sign keys.

```
vault secrets enable ssh
```

## Enable CA

Generate CA for ssh.

```
vault write -f ssh/config/ca
Key           Value
---           -----
public_key    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCwjLDgrISTB9no8Y4KE+1A7nvue1hfodlhC88jeXLbmop+mbvsG237FAs1TACHRxmyMfeM1b7N6FmLgisnf0Btk8tGPvikNcbHMf4kq7FwV2WeoSgHvtXex/dalrAnLbumthJUtYxY8dFticpxBp4guicretGbQ5LjVpOAPW8PhifrmoUgGUJ44eWET1Yg6sEiV4ZXv3nlXb7YFTJSaP3StJ4rCTFwCZA2NBV2ecBGee28+MtRx+bERecMic4GeXHXRhtHhf0gf8PaO7ePQ5pXPC1dxLG/GHAhwl04JUfc49NsWPavVuY5jUNFqfMZ4r0rN/TLl3XDyLkPPOGINkyoUoAW1Ji6DScnahYFA4i6oxQ0g3He5DWd71i6jm/3jodZ2lJjZyO8m89ceGr+f3Gwl8KVP2FQ+//psMADEcROJbK0v8YpSzJxM4tO7nsQ+FDKCbXGHAfdqwPS3TdI2AtU40FDEyf+cb41FtDCALBNdf4OoPg3gp9K2TVYvLG0iiz/1Poupqcnip04icvxFaTIAd8yzQa+9d3qzv7g6WxDpvcE04Z5VI/Lzw2wNXXOM/TdAm9YJco1vMrlJQAb8evstpSSSriMlQfeJUKFOKldgLF0cMLhgQFi5vF2R6t517N6rzy3Nxm4vKB/BomX8+iO9jXwscFSCDJt1HduXAtaQQ==
```
## Export CA public key

```
vault read -field=public_key ssh/config/ca > /tmp/ssh-ca.pem
```

This certificate will be used on ssh servers.

## Signmykey configuration

### Vault sign user role

```
echo '{                                            "allow_user_certificates": true,
  "allowed_users": "*",
  "allow_user_key_ids": true,
  "default_extensions": [
    {
      "permit-pty": ""
    }
  ],
  "key_type": "ca",
  "default_user": "root",
  "max_ttl": "24h" ,
  "ttl": "30m"
}' > sign-user-role.json
```

Then

```
vault write ssh/roles/sign-user-role @sign-user-role.json
```

### Vault policy

```
vault write sys/policy/signmykey-server policy=-<<"EOH"
path "ssh/config/ca" {
  capabilities = ["read"]
}
path "ssh/sign/sign-user-role" {
  capabilities = ["create", "update"]
}
EOH
```

### Vault AppRole

```
vault auth enable approle
```

```
vault write auth/approle/role/signmykey-server \
  token_num_uses=0 \
  token_ttl=1m \
  token_max_ttl=1m \
  policies=signmykey-server
```

### Vault AppRole creds

```
vault read auth/approle/role/signmykey-server/role-id
Key        Value
---        -----
role_id    140f639f-3c86-4bce-6019-8a9cfd4e47e8
```

``` 
vault write -f auth/approle/role/signmykey-server/secret-id
Key                   Value
---                   -----
secret_id             3f69bdc0-f854-831e-dcf7-bd2de5b8141b
secret_id_accessor    5ee59b8e-685b-f3f5-49d5-b86bf53d4424
```

The role_id and secret_id will be used in signmykey server configuration.
