package local

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestSigner(t *testing.T) {

	testCACert := []byte(`ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCxZVABBnfTplFb+U2ql7wBlH07c8+xBppg+tZwdcfy+Ib6Tj7lJRGhBmcBviPDjMeW73dntcsT6ELaRLtzP/Bo6lZBdwlVhBtDZpamOiv1qSd0L8YtYNfAuM71oyobqbtjSrEZzdhlB7nBk2usJtkol2VNqM7nnL8R99fy9uceAzZTOXutxCzu7obGVojJ0mwSgplOwcOu98BLdReUdX9j4FC6bcpTuvDdHww1e+2DU1FpY3PJMgob/DkhpUJbpClZhfroI1pgezBEox9wAUvCBbi/QBO2EV6pE4orIlGUMmkyeqKWSt6t+sqT8xEUbcG/LrJ0ZyPwmXMhk5tBjE4z`)
	testCAKey := []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAsWVQAQZ306ZRW/lNqpe8AZR9O3PPsQaaYPrWcHXH8viG+k4+
5SURoQZnAb4jw4zHlu93Z7XLE+hC2kS7cz/waOpWQXcJVYQbQ2aWpjor9akndC/G
LWDXwLjO9aMqG6m7Y0qxGc3YZQe5wZNrrCbZKJdlTajO55y/EffX8vbnHgM2Uzl7
rcQs7u6GxlaIydJsEoKZTsHDrvfAS3UXlHV/Y+BQum3KU7rw3R8MNXvtg1NRaWNz
yTIKG/w5IaVCW6QpWYX66CNaYHswRKMfcAFLwgW4v0ATthFeqROKKyJRlDJpMnqi
lkrerfrKk/MRFG3Bvy6ydGcj8JlzIZObQYxOMwIDAQABAoIBABEw3zrqNIyHLpU6
KKOihq6khCpRw8vE9wr04/kMAO9z1CjHkmLEX9v58duCYJbfuqvO0wRy4pYwSOI8
DOpTROn145v+fCIUZkv20hyTwJTS6qbgxlS5cM8VWcEGKdt6bFVn9JeqkhDgWcj/
j0ykiyDa4w9Oj5Z0YzPLj9rUwHrw7W6R3xNh69Sv6+XmX1TqbFo0Zl9nkEDfhNSh
Fi270+5HBZIfyXsFbuC6m56jtmgpDCQP51HtWiWMbXy60YPJe27HhWtOfECUtIUx
S13qPNfa3FNjfgpEHuFyGgl2qLkZLMyNe0gXR+xyArQiOhkPw4TlMUaE9da+K0WQ
ymZjl8ECgYEA2r8patPctG4jWEdZo7BOZANCP7TyPcouoDbmB88SbthOLZ71Ru9s
Xm+KxJdXaprsOK02QosAR0jdBAdKHTLsEW5FXL1FZgF3yJbTecPmql13EZVfQpxo
mK7iXpoHl0CpPWZ1S7lfa2QRurZRPfsVpMa8UvRHN9ss+JdTK1ipeSECgYEAz5tY
rQltm/M9ZIRMNoLuwgSpAhlE4Avk6O6hu3Ph+Dwf3ktsb8TuHhdCRyoUqw9xhj/I
SNSPCD15Im8IOHFGVWta5BweLxkesJnySj5wsvawGhLJ53szG0mubg2hRT1Uc5j6
lWYbibYZ80f4E2EuAe/puE/do/F32FlQhYxWeNMCgYEAlc1N/tOiFIpMeDs8nwWx
WXqF1v0C29/m+F9APt7HP9OwDjwKuw5hx3ZZsPH3spDv7oxoWT+57BdxDD41ujNS
SUmcBLu1l/qvXlYz8vJ+t/MUBJ2nxAU6+Dzj12dihWmJvPu6niYPu4qnPZd3oZue
od5bv+98Cjt127Q+B7RLMyECgYAIiLQ+eLK+xGLzrNSNMRirdRGVeoBwTUzdnmGb
mQni8GXG94a/mXLIXeBlmH89AOeDwz9ybvpqNkyyc6n81/syK7WSxu0etoOictGY
57QuRyG6EKeoElJpfr/i2kCU3g6IqfMzDP14zbmHXJ//+/CuN7R91RqhUJ3CkPlU
ZA5x7QKBgA7w9qB6RkBQUEvAghB04pwyj9u0Fj+ieNh9FjZI7ysuMYEOUIGN+3mf
aTf+DpVUfDDy7z+wU4lB5XKPp3tpVbnUvF13dJJ8qHNl79BNaso158af4mfX5A8K
i9Jm7L090D/cW6+hyuIzXBEyp+u1FXzaaEEMw9dsFyp9URYd1b03
-----END RSA PRIVATE KEY-----
`)
	testKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDXtOZp9B+Qh/DvRC75EaIoTi79TDHPQwfzIb7hcZl/Aj3NWykk8NBAIx3Mpx39sZJNtDW/AGWI9yEliSk34y4EYq3sMfgdytkVWxh3t8U6Lb7t7+I4OejCBti1tQT3g3HmxTiq1vrz7FEtZ2Ji7gc1ZaRTXvjFcayTR/7a4hSYqgnkmGUWBD3OQEP6nuC6KuU9Aur0CvVKWclRGAhE6fJI54R3D6KucL2mCWWQDlJarVv8cKDj/WVIZqt8CdgE46HNxCRcSDUzpTUSqThPM9oqaPC6xBLjpBEwTQpvOnYufB4TOepjdy+221cLkaflgn0JZZyU/39VNXbI7no/VPwh"

	CACert, _, _, _, err := ssh.ParseAuthorizedKey(testCACert)
	if err != nil {
		t.Error(err)
	}

	CAKey, err := ssh.ParsePrivateKey(testCAKey)
	if err != nil {
		t.Error(err)
	}

	s := &Signer{
		CACert: CACert,
		CAKey:  CAKey,
		TTL:    600,
	}

	cases := []struct {
		description string
		payload     []byte
		id          string
		principals  []string
		expErr      bool
	}{
		{"test with an invalid key", []byte("{\"public_key\": \"invalid key\"}"), "test", []string{"root", "admin"}, true},
		{"test with valid key and principals", []byte(fmt.Sprintf("{\"public_key\": \"%s\"}", testKey)), "testid", []string{"admin", "root"}, false},
		{"test with valid key and reversed principals", []byte(fmt.Sprintf("{\"public_key\": \"%s\"}", testKey)), "testid", []string{"root", "admin"}, false},
		{"test with an empty key", []byte(""), "testid", []string{"root", "admin"}, true},
		{"test with an empty id", []byte(fmt.Sprintf("{\"public_key\": \"%s\"}", testKey)), "", []string{"root", "admin"}, true},
		{"test with no principals", []byte(fmt.Sprintf("{\"public_key\": \"%s\"}", testKey)), "testid", []string{}, true},
	}

	for _, c := range cases {
		cert, err := s.Sign(context.Background(), c.payload, c.id, c.principals)
		if c.expErr {
			assert.Error(t, err, c.description)
		} else {
			assert.NoError(t, err, c.description)
		}

		if err != nil {
			continue
		}

		parsedCert, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(cert))
		sshCert := parsedCert.(*ssh.Certificate)

		sort.Strings(c.principals)
		sort.Strings(sshCert.ValidPrincipals)

		assert.Equal(t, c.id, sshCert.KeyId, c.description)
		assert.Equal(t, c.principals, sshCert.ValidPrincipals, c.description)
	}
}
