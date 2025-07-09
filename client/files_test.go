package client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestChooseSSHKeyType(t *testing.T) {
	cases := []struct {
		keyName       string
		keyType       string
		keyDeprecated bool
	}{
		{"~/.ssh/id_dsa.pub", "dsa", true},
		{"~/.ssh/id_ecdsa.pub", "ecdsa", false},
		{"~/.ssh/id_ecdsa_sk.pub", "ecdsa-sk", false},
		{"~/.ssh/id_ed25519.pub", "ed25519", false},
		{"~/.ssh/id_ed25519_sk.pub", "ed25519-sk", false},
		{"~/.ssh/id_rsa.pub", "rsa", false},
		{"~/.ssh/test_default_type.pub", "ed25519", false},
	}

	for _, c := range cases {
		keyType, isDeprecated := chooseSSHKeyType(c.keyName)
		assert.Equal(t, c.keyType, keyType)
		assert.Equal(t, c.keyDeprecated, isDeprecated)
	}
}

func TestCertKeyTypeIsBuggy(t *testing.T) {
	cases := []struct {
		keyType string
		isBuggy bool
	}{
		{ssh.KeyAlgoRSA, true},
		{ssh.KeyAlgoDSA, false}, //nolint:staticcheck
		{ssh.KeyAlgoED25519, false},
		{ssh.KeyAlgoECDSA256, false},
	}

	for _, c := range cases {
		isBuggy := CertKeyTypeIsBuggy(c.keyType)
		assert.Equal(t, c.isBuggy, isBuggy, c.keyType)
	}
}

func TestCertInfo(t *testing.T) {
	cases := []struct {
		principals []string
		before     uint64
		keyType    string
		cert       string
	}{
		{[]string{"test_principal"}, 18446744073709551615, ssh.KeyAlgoRSA, "ssh-rsa-cert-v01@openssh.com AAAAHHNzaC1yc2EtY2VydC12MDFAb3BlbnNzaC5jb20AAAAg7y+wbqFwFYKHvB8iBG3fSUiFDJwlwrgYK29SVcfXIl8AAAADAQABAAABgQDXyccdoqiuomunLwK/u7cddjovr38Svi5M7qq3SoUSk0Y5BakECmB2QlevqXht6o938kia3THDv7HGMoEIwhuzAR/2pEqobr7bTo6yfxuBjiGU2Tbu4Ma7/gDfWlMavFKyTlOiIeQQXn+I6VG4UW9ddLc76u3VhStclW6RjdMQ1RAeO4ln+/GYTkoa8SjLfw75gZe1tXDlMxQrPK7vvoqyY7a+NYomrnatozI7mKUCxl1Du1Rkp8o23JSSGNxRi2Q7iNHE/3I6XJ8/g67Yrniky0NEyXtP9hB5u74hnnr9eAOCC/d1LKiVvlwHC8Vc8/e3DVDE4tcEvI6njIi9w4YOZ6/9aXXEzyUunGPTh0M09paEu1dgb5bwVnPs0/50g9N4efs/s0QRVvW5osAxkdjoT9krIK/pgHEZRlRN+4T1Ogj4erl2ZMqb1VTmeCSL69Ck6upotbkTUpmBuIDXK1HGt83g/HEc7LC1OMPm3UoDGndYuRORu54qJg2UHaHz9IcAAAAAAAAAAAAAAAEAAAAHdGVzdF9pZAAAABIAAAAOdGVzdF9wcmluY2lwYWwAAAAAAAAAAP//////////AAAAAAAAAIIAAAAVcGVybWl0LVgxMS1mb3J3YXJkaW5nAAAAAAAAABdwZXJtaXQtYWdlbnQtZm9yd2FyZGluZwAAAAAAAAAWcGVybWl0LXBvcnQtZm9yd2FyZGluZwAAAAAAAAAKcGVybWl0LXB0eQAAAAAAAAAOcGVybWl0LXVzZXItcmMAAAAAAAAAAAAAAZcAAAAHc3NoLXJzYQAAAAMBAAEAAAGBALLjeJw63Jc5dNQFmH50B7zRGfYzRErgJeyDzJfPO3LKHcqpOt5TQhkkO8pl2KKZmfoUoZWr158B8hdh/5khexWSVRGacb4K+cdxbO1zA6LXtH40KD2wlfYHTrA3ze4NCJ4jKrfJE8uNjnukp78mYtsuwvmwJLsRcXHXxLpMmSRsdSkf37k0BJ+maHkrXzrHhpLfu+YcNW56MKUer8+Bg3GruUPQNlhtJlIvnMJ702k8DmHV4Gdcal+ixgCDPAoZMzEYuqILtsb/gOVQPtozM/kILyVTJlBA8y1/wIjgQkcpy7vOaIOrvNjEu1FKrz/+lJEOqC4YNQs+ad2NCfvJPs8S5Qe/J/bdukGLRITmZ9FeQOyvgdf7lp+UWNAN1oiO4406DoEBV8semFwQuGyz8Co4jy1MWIBRPNrRRcB3Lkl3lpPz3uj3UQRguqpEKXp2j8gcdyhF/Ci9x+k8WYf97U+Szxsj71ZnJ/buKShFo87S44KoIAUrq5tFxI1qbuFTxQAAAZQAAAAMcnNhLXNoYTItNTEyAAABgHLb4k831Pv9RmDcxgQfZ172GVXGYuWFJiAg2AiCWy4Va5Pp4IaNFFwEQ/KtLDLfb7z18vQ+cZN1Yd6eWjLfLeqoo1K4ma0hmofNTvZ2kp70bYSDBB5mSSFifigbPrn1b6B7vJhmKZbUITesW1pKwDYOPeLv71FbLf1Z+nkVyhHj4aGXEBiZbO+FvXmWiJYs1yJHio9c+SKBq2fr3g0LBhcYAqxroCWVbji7SfzyU96vEqve0nscWQ8HSFxVAXAObc5/YBnCBQrxHzIbrpD94Wr2y9r7SOLCnA2ttEw9IUf8Vgt/s122aoLU13lOpjJb0Or4meW2pxBOMbAH4NCbAfqkDRIPzZz8RDg8kw77g8q/Mbr7on+WyjOhJoCPeZgGwVC9JRm0sXNpiziJT7Nl3rIChNkAf6ysyRF0TOHpnEEsV5u5z3o4kftjHdBh3xDKSRlQ7R+UgtM/bjD+gy9ImwuqbO4ZE0SlpgJsXkUV2N3TA1vAxbhjFakYZqvJmZ786w=="},
		{[]string{"test_principal"}, 18446744073709551615, ssh.KeyAlgoECDSA256, "ecdsa-sha2-nistp256-cert-v01@openssh.com AAAAKGVjZHNhLXNoYTItbmlzdHAyNTYtY2VydC12MDFAb3BlbnNzaC5jb20AAAAgJdvGTH0BYisEPsXsuF0YzbiQEzTbu8IcS84/hhboUS0AAAAIbmlzdHAyNTYAAABBBNVjZJiuR4CNpii/d16CF3mst0pvNFVC3y5iIqhRo/p86JFF3nsKliSRYUvTA+kFGJbMr7uL41KH/qrmiQL7NsoAAAAAAAAAAAAAAAEAAAAHdGVzdF9pZAAAABIAAAAOdGVzdF9wcmluY2lwYWwAAAAAAAAAAP//////////AAAAAAAAAIIAAAAVcGVybWl0LVgxMS1mb3J3YXJkaW5nAAAAAAAAABdwZXJtaXQtYWdlbnQtZm9yd2FyZGluZwAAAAAAAAAWcGVybWl0LXBvcnQtZm9yd2FyZGluZwAAAAAAAAAKcGVybWl0LXB0eQAAAAAAAAAOcGVybWl0LXVzZXItcmMAAAAAAAAAAAAAAZcAAAAHc3NoLXJzYQAAAAMBAAEAAAGBALLjeJw63Jc5dNQFmH50B7zRGfYzRErgJeyDzJfPO3LKHcqpOt5TQhkkO8pl2KKZmfoUoZWr158B8hdh/5khexWSVRGacb4K+cdxbO1zA6LXtH40KD2wlfYHTrA3ze4NCJ4jKrfJE8uNjnukp78mYtsuwvmwJLsRcXHXxLpMmSRsdSkf37k0BJ+maHkrXzrHhpLfu+YcNW56MKUer8+Bg3GruUPQNlhtJlIvnMJ702k8DmHV4Gdcal+ixgCDPAoZMzEYuqILtsb/gOVQPtozM/kILyVTJlBA8y1/wIjgQkcpy7vOaIOrvNjEu1FKrz/+lJEOqC4YNQs+ad2NCfvJPs8S5Qe/J/bdukGLRITmZ9FeQOyvgdf7lp+UWNAN1oiO4406DoEBV8semFwQuGyz8Co4jy1MWIBRPNrRRcB3Lkl3lpPz3uj3UQRguqpEKXp2j8gcdyhF/Ci9x+k8WYf97U+Szxsj71ZnJ/buKShFo87S44KoIAUrq5tFxI1qbuFTxQAAAZQAAAAMcnNhLXNoYTItNTEyAAABgJsWz2/FgCf/mfnkw2Yn7BFFcH3ZAYsO89KYTJmCzRit0o5EcHjoFtPBf7jeiieJ19nnDR4jpC+5FwqNnXiI0xpSus4Rh3D84h78xxHIRhR3nBqfK0CCP4SXl7aZLIvOf32PNLvVXehSXMfq04hwzEqYOdHBf37O5uwdukFKX0FSEzGiEPa51LPEzs8wugwNzp/IKnO/DOI8D8i5cySxGMp8Rfai1dRquPd0NxBq4qIDaDJ9wz31YGOEEAxIdWQCvwGx0hUU7U6DD/8W2NMEtyVIPZPh7dzZOgXWpIuWQ08v/yqabXZEMqayjR3yQ9/cVnFix9SdjOo3xaP9O/O5lX6jWWHWkEhU80QnlJQRgcPIoYP4LXeNVvi8dDEFOLd6m+bO3/ysjm4z+GzQhxzOuCV7Xh8BeoB9EDG3LGgROIzRWp6NWk1GEOa4S8bu/coo/fQkXuVn61SgB2Km2RRlqB3R4CSNgX5/cawczvyGs4fBcAIYmILDyg4qbvbDSKLvPg=="},
		{[]string{"test_principal"}, 18446744073709551615, ssh.KeyAlgoED25519, "ssh-ed25519-cert-v01@openssh.com AAAAIHNzaC1lZDI1NTE5LWNlcnQtdjAxQG9wZW5zc2guY29tAAAAICenegcHqkkVK5eZApDvOKoqK6TSv4yAdjyGcFNreh6bAAAAIDqolxIZDb/1QnCQfVnlv4Uy4IMjfWMmUYIn6DXG1askAAAAAAAAAAAAAAABAAAAB3Rlc3RfaWQAAAASAAAADnRlc3RfcHJpbmNpcGFsAAAAAAAAAAD//////////wAAAAAAAACCAAAAFXBlcm1pdC1YMTEtZm9yd2FyZGluZwAAAAAAAAAXcGVybWl0LWFnZW50LWZvcndhcmRpbmcAAAAAAAAAFnBlcm1pdC1wb3J0LWZvcndhcmRpbmcAAAAAAAAACnBlcm1pdC1wdHkAAAAAAAAADnBlcm1pdC11c2VyLXJjAAAAAAAAAAAAAAGXAAAAB3NzaC1yc2EAAAADAQABAAABgQCy43icOtyXOXTUBZh+dAe80Rn2M0RK4CXsg8yXzztyyh3KqTreU0IZJDvKZdiimZn6FKGVq9efAfIXYf+ZIXsVklURmnG+CvnHcWztcwOi17R+NCg9sJX2B06wN83uDQieIyq3yRPLjY57pKe/JmLbLsL5sCS7EXFx18S6TJkkbHUpH9+5NASfpmh5K186x4aS37vmHDVuejClHq/PgYNxq7lD0DZYbSZSL5zCe9NpPA5h1eBnXGpfosYAgzwKGTMxGLqiC7bG/4DlUD7aMzP5CC8lUyZQQPMtf8CI4EJHKcu7zmiDq7zYxLtRSq8//pSRDqguGDULPmndjQn7yT7PEuUHvyf23bpBi0SE5mfRXkDsr4HX+5aflFjQDdaIjuONOg6BAVfLHphcELhss/AqOI8tTFiAUTza0UXAdy5Jd5aT897o91EEYLqqRCl6do/IHHcoRfwovcfpPFmH/e1Pks8bI+9WZyf27ikoRaPO0uOCqCAFK6ubRcSNam7hU8UAAAGUAAAADHJzYS1zaGEyLTUxMgAAAYAM8l4R33l3uJ2nRjm8y5e4b98mkJ6q2tta7N3PbrOB5jnMkYv7TKDqPBnawWkP2w6+LA1M/7XWvPmNro5zkUeeSs9UUrXS9OVKN69H/cM31jNLcXeaUOE0BwSC6jXWLIsPlcC8oNXAble7fGfdzCVqgno0jU/7E/I7IMfk4K0Fo8sFdNVoT0gazVKQkY1nb7WpSgdEtMFlns2wQ8mwZBbgVUE2w7AYz78CTdozajx4EXQffQjsrPbYXA3F9Icgc7CnDyH0O5nOb+YSYBnAXFSf7UCPIzBDA4n6jp+noGwHu2EGfpGGKD+KZx0KYStIGnXGYYF+1wWK2vo3xe8U1oQ9N/EX0dvt0Pwi6WsIIpemqAvjCWz/aQwedFD8K8hEK84OBzvUxv0cAm2aWalBb+PYJl8mYmkrsgWrcb/XX8VCBka3bgKK2lw8023zFe9NvdvSKx/a6/qXspxg/zMRwaUUMzre5hmG/uRB2KoDTjIrwlpYseOGUrTAFT0jH/97Ikc="},
	}

	for _, c := range cases {
		principals, before, keyType, _ := CertInfo(c.cert)
		if !reflect.DeepEqual(principals, c.principals) {
			t.Errorf("CertInfo principal: %v, want %v", principals, c.principals)
		}

		if before != c.before {
			t.Errorf("CertInfo before: %v, want %v", before, c.before)
		}

		if keyType != c.keyType {
			t.Errorf("CertInfo keyType: %v, want %v", keyType, c.keyType)
		}
	}
}
