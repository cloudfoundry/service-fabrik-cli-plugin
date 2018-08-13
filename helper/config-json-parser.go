package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-incubator/service-fabrik-cli-plugin/errors"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	RefreshToken          string
	AccessToken           string
	SpaceFields           SpaceField
	OrganizationFields    OrgField
	Target                string
	AuthorizationEndpoint string
}

type SpaceField struct {
	GUID string
	Name string
}

type OrgField struct {
	GUID string
	Name string
}

type TokenInfo struct {
	Username string   `json:"user_name"`
	Email    string   `json:"email"`
	UserGUID string   `json:"user_id"`
	GUID     string   `json:"GUID"`
	Scope    []string //`json:"scope":["cloud_controller.read","password.write","cloud_controller.write","openid","uaa.user"]`
}

//this code taken from cf cli source code source
func NewTokenInfo(accessToken string) (info TokenInfo) {
	tokenJSON, err := DecodeAccessToken(accessToken)
	if err != nil {
		return TokenInfo{}
	}

	info = TokenInfo{}
	err = json.Unmarshal(tokenJSON, &info)
	if err != nil {
		return TokenInfo{}
	}

	return info
}

func DecodeAccessToken(accessToken string) (tokenJSON []byte, err error) {
	tokenParts := strings.Split(accessToken, " ")

	if len(tokenParts) < 2 {
		return
	}

	token := tokenParts[1]
	encodedParts := strings.Split(token, ".")

	if len(encodedParts) < 3 {
		return
	}

	encodedTokenJSON := encodedParts[1]
	return base64Decode(encodedTokenJSON)
}

func base64Decode(encodedData string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(restorePadding(encodedData))
}

func restorePadding(seg string) string {
	switch len(seg) % 4 {
	case 2:
		seg = seg + "=="
	case 3:
		seg = seg + "="
	}
	return seg
}

type Test1 struct {
	Msg    string
	Field1 Field1Type
	Fields []FieldsType
}

type Field1Type struct {
	Msg string
}

type FieldsType struct {
	Msg string
}

func ReadConfigJsonFile() []byte {
	var CF_HOME string = os.Getenv("CF_HOME")

	if CF_HOME == "" {
		CF_HOME = GetHomeDir()
	}

	file, err := ioutil.ReadFile(CF_HOME + string(os.PathSeparator) + ".cf" + string(os.PathSeparator) + "config.json")
	if err != nil {
		errors.FileReadingError("config.json")
	}
	return file
}

func GetHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		errors.HomeDirNotFound(err)
	}
	return homeDir
}

func GetAccessToken(file []byte) string {
	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Error parsing json [%v]\n", err)
		os.Exit(4)
	}

	if config.AccessToken == "" {
		errors.NoAccessTokenError("Access Token")
	}
	return config.AccessToken
}

func GetRefreshToken(file []byte) string {
	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Error parsing json [%v]\n", err)
		os.Exit(4)
	}

	if config.AccessToken == "" {
		errors.NoAccessTokenError("Access Token")
	}
	return config.RefreshToken
}

func GetSpaceGUID(file []byte) string {
	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Error parsing json [%v]\n", err)
		os.Exit(4)
	}

	if config.SpaceFields.GUID == "" {
		errors.NoAccessTokenError("Space Fields")
	}
	var userSpaceGuid string = strings.Trim(config.SpaceFields.GUID, "\"")
	return userSpaceGuid
}

func GetSpaceName(file []byte) string {
	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Error parsing json [%v]\n", err)
		os.Exit(4)
	}

	if config.SpaceFields.Name == "" {
		errors.NoAccessTokenError("Space Fields")
	}
	var userSpaceName string = strings.Trim(config.SpaceFields.Name, "\"")
	return userSpaceName
}

func GetOrgName(file []byte) string {
	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Error parsing json [%v]\n", err)
		os.Exit(4)
	}

	if config.OrganizationFields.Name == "" {
		errors.NoAccessTokenError("Organisation Fields")
	}
	var userOrgName string = strings.Trim(config.OrganizationFields.Name, "\"")
	return userOrgName
}

func GetApiEndpoint(file []byte) string {
	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Error parsing json [%v]\n", err)
		os.Exit(4)
	}

	if config.Target == "" {
		errors.NoAccessTokenError("Api Endpoint")
	}
	return config.Target

}

func GetLoginEndpoint(file []byte) string {
	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Error parsing json [%v]\n", err)
		os.Exit(4)
	}

	if config.Target == "" {
		errors.NoAccessTokenError("Login Endpoint")
	}
	return config.AuthorizationEndpoint

}

func Exists(path string) bool {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateConfFile() {
	var path string
	var CF_HOME string = os.Getenv("CF_HOME")
	if CF_HOME == "" {
		CF_HOME = GetHomeDir()
	}
	path = CF_HOME + "/.cf/conf.json"

	brace1 := []byte("{\n")
	key1 := []byte("\"serviceBroker\": ")
	val1 := []byte("\"service-fabrik-broker\",\n")
	key2 := []byte("\"serviceBrokerExtUrl\": ")
	val2 := []byte("\"/api/v1\",\n")
	key3 := []byte("\"skipSslFlag\": ")
	val3 := []byte("true\n")
	brace2 := []byte("}")
	if Exists(path) {
		return
	} else {
		f, err := os.Create(path)
		errors.ErrorIsNil(err)

		defer f.Close()

		f.Write(brace1)
		f.Write(key1)
		f.Write(val1)
		f.Write(key2)
		f.Write(val2)
		f.Write(key3)
		f.Write(val3)
		f.Write(brace2)

		f.Sync()
	}
}
