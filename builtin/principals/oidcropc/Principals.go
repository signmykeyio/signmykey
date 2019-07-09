package oidcropc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/signmykeyio/signmykey/builtin/principals/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Principals struct represents OIDC options for getting principals list from OIDC.
type Principals struct {
	OIDCUserinfoEndpoint string
	OIDCUserGroupsEntry  string
	TransformCase        string
}

type oidcUserinfo struct {
	Oidcgroups []string `json:"oidcgroups"`
}

// Init method is used to ingest config of Principals
func (p *Principals) Init(config *viper.Viper) error {
	neededEntries := []string{
		"oidcUserinfoEndpoint",
		"oidcUserGroupsEntry",
	}

	var missingEntriesLst []string
	for _, entry := range neededEntries {
		if !config.IsSet(entry) {
			missingEntriesLst = append(missingEntriesLst, entry)
			continue
		}
	}
	if len(missingEntriesLst) > 0 {
		missingEntries := strings.Join(missingEntriesLst, ", ")
		return fmt.Errorf("Missing config entries (%s) for Principals", missingEntries)
	}

	config.SetDefault("transformCase", "none")
	tc := config.GetString("transformCase")
	if tc != "none" && tc != "lower" && tc != "upper" {
		return errors.New("transformCase config entry for Principals must be none, lower or upper")
	}

	p.OIDCUserinfoEndpoint = config.GetString("oidcUserinfoEndpoint")
	p.OIDCUserGroupsEntry = config.GetString("oidcUserGroupsEntry")

	p.TransformCase = tc

	return nil
}

// Get method is used to get the list of principals associated to a specific user.
func (p Principals) Get(user string) (principals []string, err error) {

	reqInfo, err := http.NewRequest("GET", p.OIDCUserinfoEndpoint, nil)
	if err != nil {
		return principals, err
	}

	// Add HTTP Authorization Header
	bearer := "Bearer " + user
	reqInfo.Header.Add("Authorization", bearer)

	client := http.Client{Timeout: time.Second * 10}
	resInfo, err := client.Do(reqInfo)
	if err != nil {
		return principals, err
	}

	defer resInfo.Body.Close()

	bodyInfo, err := ioutil.ReadAll(resInfo.Body)
	if err != nil {
		return principals, err
	}

	// Replace `json:"oidcgroups"` oidcUserinfo struct tag with OIDCUserGroupsEntry config entry
	bodyInfoChange := []byte(strings.Replace(string(bodyInfo), p.OIDCUserGroupsEntry, "oidcgroups", 1))

	oidcUserinfo1 := oidcUserinfo{}
	err = json.Unmarshal(bodyInfoChange, &oidcUserinfo1)
	if err != nil {
		return principals, err
	}

	principals = oidcUserinfo1.Oidcgroups
	logrus.Debug("OIDC principals: ", principals)

	principals = common.TransformCase(p.TransformCase, principals)

	return principals, nil
}
