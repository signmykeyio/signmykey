---
title: Vault
---

Before you can use Signmykey, you must configure your Vault server to allow its usage.

## Enable ssh

First, enable the ssh secret engine to sign keys.

```sh
vault secrets enable ssh
```

## Enable CA

Generate CA for ssh.

```sh
vault write -f ssh/config/ca
Key           Value
---           -----
public_key    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCwjLDgrISTB9no8Y4KE+1A7nvue1hfodlhC88jeXLbmop+mbvsG237FAs1TACHRxmyMfeM1b7N6FmLgisnf0Btk8tGPvikNcbHMf4kq7FwV2WeoSgHvtXex/dalrAnLbumthJUtYxY8dFticpxBp4guicretGbQ5LjVpOAPW8PhifrmoUgGUJ44eWET1Yg6sEiV4ZXv3nlXb7YFTJSaP3StJ4rCTFwCZA2NBV2ecBGee28+MtRx+bERecMic4GeXHXRhtHhf0gf8PaO7ePQ5pXPC1dxLG/GHAhwl04JUfc49NsWPavVuY5jUNFqfMZ4r0rN/TLl3XDyLkPPOGINkyoUoAW1Ji6DScnahYFA4i6oxQ0g3He5DWd71i6jm/3jodZ2lJjZyO8m89ceGr+f3Gwl8KVP2FQ+//psMADEcROJbK0v8YpSzJxM4tO7nsQ+FDKCbXGHAfdqwPS3TdI2AtU40FDEyf+cb41FtDCALBNdf4OoPg3gp9K2TVYvLG0iiz/1Poupqcnip04icvxFaTIAd8yzQa+9d3qzv7g6WxDpvcE04Z5VI/Lzw2wNXXOM/TdAm9YJco1vMrlJQAb8evstpSSSriMlQfeJUKFOKldgLF0cMLhgQFi5vF2R6t517N6rzy3Nxm4vKB/BomX8+iO9jXwscFSCDJt1HduXAtaQQ==
```
## Export CA public key

```sh
vault read -field=public_key ssh/config/ca > /etc/ssh/ca.pem
```

This certificate will be used on ssh servers, so keep it to copy its content later.

## Signmykey configuration

### Vault sign user role

```sh
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

```sh
vault write ssh/roles/sign-user-role @sign-user-role.json
```

### Vault policy

```sh
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

```sh
vault auth enable approle
```

```sh
vault write auth/approle/role/signmykey-server \
  token_num_uses=0 \
  token_ttl=1m \
  token_max_ttl=1m \
  policies=signmykey-server
```

### Vault AppRole creds

```sh
vault read auth/approle/role/signmykey-server/role-id
Key        Value
---        -----
role_id    11940c2d-4639-9358-d750-cdb7cf409ff4
```

```sh
vault write -f auth/approle/role/signmykey-server/secret-id
Key                   Value
---                   -----
secret_id             8b4c901f-1f84-5049-17ee-92de12b6b1e5
secret_id_accessor    0921e287-5383-0fbd-5061-aef29618b7a0
```

The role_id and secret_id will be used in signmykey server configuration.
