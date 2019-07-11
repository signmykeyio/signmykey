package local

import (
	"bytes"
	"context"
	"sort"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPrincipals(t *testing.T) {

	cases := []struct {
		userMap []byte
		payload []byte
		expErr  bool
		expList []string
	}{
		{
			[]byte(`
users:
  user1: princ2,princ1`),
			[]byte("{\"user\": \"user1\"}"), false, []string{"princ1", "princ2"},
		},
		{
			[]byte(`
users:
  user: princ1,princ2`),
			[]byte("{\"user\": \"user\"}"), false, []string{"princ1", "princ2"},
		},
		{
			[]byte(`
users:
  user2: princ1,princ2`),
			[]byte("{\"user\": \"user1\"}"), true, []string{},
		},
		{
			[]byte(`
users:
  user1: princ1
  user2: princ3,princ4
`),
			[]byte("{\"user\": \"user2\"}"), false, []string{"princ3", "princ4"},
		},
		{
			[]byte(`
users:
  user1: princ1, princ2,princ3
  user2: princ3,princ4
`),
			[]byte("{\"user\": \"user1\"}"), false, []string{"princ1", "princ2", "princ3"},
		},
		{
			[]byte(`
users:
  user1: princ1, princ2,princ3 , ,princ4
  user2: princ3,princ4
`),
			[]byte("{\"user\": \"user1\"}"), false, []string{"princ1", "princ2", "princ3", "princ4"},
		},
	}

	for _, c := range cases {
		testConfig := viper.New()
		testConfig.SetConfigType("yaml")
		err := testConfig.ReadConfig(bytes.NewBuffer(c.userMap))
		if err != nil {
			t.Error(err)
		}

		local := &Principals{}
		local.Init(testConfig)

		_, principals, err := local.Get(context.Background(), c.payload)
		if c.expErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		sort.Strings(principals)
		sort.Strings(c.expList)
		assert.Equal(t, c.expList, principals)
	}
}
