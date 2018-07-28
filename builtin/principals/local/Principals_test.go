package local

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrincipals(t *testing.T) {

	cases := []struct {
		userMap map[string]string
		user    string
		expErr  bool
		expList []string
	}{
		{
			map[string]string{"user1": "princ2,princ1"}, "user1",
			false, []string{"princ1", "princ2"},
		},
		{
			map[string]string{"user": "princ1,princ2"}, "user",
			false, []string{"princ1", "princ2"},
		},
		{
			map[string]string{"user2": "princ1,princ2"}, "user1",
			true, []string{},
		},
		{
			map[string]string{"user1": "princ1", "user2": "princ3,princ4"}, "user2",
			false, []string{"princ3", "princ4"},
		},
		{
			map[string]string{"user1": "princ1, princ2,princ3 ", "user2": "princ3,princ4"}, "user1",
			false, []string{"princ1", "princ2", "princ3"},
		},
		{
			map[string]string{"user1": "princ1, princ2,princ3 , ,princ4", "user2": "princ3,princ4"}, "user1",
			false, []string{"princ1", "princ2", "princ3", "princ4"},
		},
		{
			map[string]string{"user1": "  , ,    , ", "user2": "princ3,princ4"}, "user1",
			true, []string{},
		},
	}

	for _, c := range cases {
		local := &Principals{}
		local.Init(c.userMap)

		principals, err := local.Get(c.user)
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
