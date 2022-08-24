package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	BaseUrl   string
	AuthUrl   string
	AuthToken string
	RoleID    string
	SecretID  string
	HttpProxy string
)

func Login() error {
	if !IsValidString(AuthUrl) {
		return errors.New("invalid auth url")
	}

	if !IsValidString(BaseUrl) {
		return errors.New("invalid base url")
	}

	if !IsValidString(RoleID) {
		return errors.New("invalid role id")
	}

	if !IsValidString(SecretID) {
		return errors.New("invalid secret id")
	}

	reqMap := map[string]string{"role_id": RoleID, "secret_id": SecretID}
	reqBody, err := json.Marshal(reqMap)
	if err != nil {
		return err
	}

	authResp, err := Post(AuthUrl, HttpProxy, reqBody)
	if err != nil {
		if err.Error() == "401" {
			return errors.New("invalid credentials")
		}
		return errors.New(fmt.Sprintf("error duing vault login: %v", err))
	}

	authResult := map[string]interface{}{}

	err = json.Unmarshal(authResp, &authResult)
	if err != nil {
		return errors.New(fmt.Sprintf("Error during parse vault: %v", err))
	}
	if v, found := authResult["auth"]; found {
		t := v.(map[string]interface{})
		if token, found := t["client_token"]; found {
			AuthToken = token.(string)
		}
	}
	return nil
}

func ReadData(dataPath string) (*map[string]interface{}, error) {
	path := formatDataPath(dataPath)
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("invalid secret path")
	}

	// add / to the end of data url if not exist
	if BaseUrl[len(BaseUrl)-1:] != "/" {
		BaseUrl += "/"
	}

	dataUrl := BaseUrl + path
	output := map[string]interface{}{}

	fmt.Printf("Vault data url: %s", dataUrl)

	// The token expire time is set by vault, can set to unlimit by set lease_duration to 0 and renewable to false
	if AuthToken == "" {
		return nil, errors.New("vault token not defined")
	}

	fmt.Printf("Read vault data token: %v\n", AuthToken)

	headers := map[string]string{"X-Vault-Token": AuthToken}
	out, err := Get(dataUrl, HttpProxy, headers)
	if err != nil {
		return nil, err
	}

	rawData := map[string]interface{}{}
	err = json.Unmarshal(out, &rawData)
	if err != nil {
		return nil, err
	}

	if v, found := rawData["data"]; found {
		if t, found := v.(map[string]interface{})["data"]; found {
			output = t.(map[string]interface{})
		}
	}

	return &output, nil
}

// remove the / prefix ///test/ -> test/
func formatDataPath(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}

	result := strings.TrimSpace(s)
	if result[0:1] == "/" {
		result = formatDataPath(result[1:])
	}

	return result
}

func IsValidString(obj string) bool {
	return obj != "" && strings.TrimSpace(obj) != ""
}
