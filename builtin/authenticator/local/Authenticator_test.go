package local

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticator(t *testing.T) {
	configBytes := []byte(`
users:
  gooduser: "$2a$10$h8bTe02uZIkAa5j1NiuVVOXdUONmch.y151qyK004Hb8EF7rTRq0u"
`)

	testConfig := viper.New()
	testConfig.SetConfigType("yaml")
	err := testConfig.ReadConfig(bytes.NewBuffer(configBytes))
	if err != nil {
		t.Error(err)
	}

	local := &Authenticator{}
	err = local.Init(testConfig)
	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		payload []byte
		id      string
		err     string
		valid   bool
	}{
		{[]byte(""), "", "JSON unmarshaling failed: unexpected end of JSON input", false},
		{[]byte("{\"user\":\"\"}"), "", "empty username", false},
		{[]byte("{\"password\":\"\"}"), "", "empty username", false},
		{[]byte("{\"user\":\"baduser\"}"), "", "empty password", false},
		{[]byte("{\"user\":\"baduser\",\"password\":\"badpassword\"}"), "", "user not found", false},
		{[]byte("{\"user\":\"gooduser\",\"password\":\"badpassword\"}"), "", "bad password", false},
		{[]byte("{\"user\":\"gooduser\",\"password\":\"goodpassword\"}"), "local-gooduser", "", true},
	}

	for _, c := range cases {
		_, valid, id, err := local.Login(context.Background(), c.payload)

		assert.Equal(t, c.valid, valid)

		assert.Equal(t, id, c.id)

		if !c.valid {
			assert.EqualError(t, err, c.err)
		}
	}
}

func FuzzAuthenticator(f *testing.F) {

	configBytes := []byte(`
users:
  gooduser: "$2a$10$h8bTe02uZIkAa5j1NiuVVOXdUONmch.y151qyK004Hb8EF7rTRq0u"
`)

	testConfig := viper.New()
	testConfig.SetConfigType("yaml")
	err := testConfig.ReadConfig(bytes.NewBuffer(configBytes))
	if err != nil {
		panic(err)
	}

	local := &Authenticator{}
	err = local.Init(testConfig)
	if err != nil {
		panic(err)
	}

	f.Fuzz(func(t *testing.T, user, password string) {

		_, _, _, err := local.Login(context.Background(), []byte(fmt.Sprintf("{\"user\":\"%s\",\"password\":\"%s\"}", user, password)))
		if err == nil && user != "gooduser" && password != "goodpassword" {
			t.Fail()
		}

	})
}

func FuzzAuthenticatorOTP(f *testing.F) {

	configBytes := []byte(`
users:
  gooduser: "$2a$10$odlh6WXIuQGMoQq5/qfJI.Q/20MGQW.NyRQVBClKMzcUaX5SElzx2,75FPOXNOPVUAVC4TQUGFD6AOFUQTUZNJ3FURT4NUJOGM5HQZELQQ===="
`)

	testConfig := viper.New()
	testConfig.SetConfigType("yaml")
	err := testConfig.ReadConfig(bytes.NewBuffer(configBytes))
	if err != nil {
		panic(err)
	}

	local := &Authenticator{}
	err = local.Init(testConfig)
	if err != nil {
		panic(err)
	}

	f.Fuzz(func(t *testing.T, user, password, otp string) {

		_, _, _, err := local.Login(context.Background(), []byte(fmt.Sprintf("{\"user\":\"%s\",\"password\":\"%s\",\"otp\": \"%s\"}", user, password, otp)))
		if err == nil && user != "gooduser" && password != "goodpassword" {
			t.Fail()
		}

	})
}
