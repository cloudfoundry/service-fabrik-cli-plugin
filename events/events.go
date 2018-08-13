package events

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-incubator/service-fabrik-cli-plugin/helper"
	"github.com/cloudfoundry-incubator/service-fabrik-cli-plugin/constants"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	httpClientDisabledSecurityCheck *http.Client
	httpClientEnabledSecurityCheck  *http.Client
)

type EventCommand struct {
	cliConnection plugin.CliConnection
}

func Initialize() {
	httpClientDisabledSecurityCheck = CreateHttpClient(true)
	httpClientEnabledSecurityCheck = CreateHttpClient(false)
}

func NewEventsCommand(cliConnection plugin.CliConnection) *EventCommand {
	command := new(EventCommand)
	command.cliConnection = cliConnection
	return command
}

func AddColor(text string, textColor color.Attribute) string {
	printer := color.New(textColor).Add(color.Bold).SprintFunc()
	return printer(text)
}

type Configuration struct {
	ServiceBroker       string
	ServiceBrokerExtUrl string
	SkipSslFlag         bool
}

func GetConfiguration() Configuration {
	var path string
	var CF_HOME string = os.Getenv("CF_HOME")
	if CF_HOME == "" {
		CF_HOME = helper.GetHomeDir()
	}
	path = CF_HOME + "/.cf/conf.json"
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}

func CreateHttpClient(disableSecurityCheck bool) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: constants.MaxIdleConnections,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: disableSecurityCheck},
			Proxy:           http.ProxyFromEnvironment,
		},
		Timeout: time.Duration(constants.RequestTimeout) * time.Second,
	}
	return client
}

func CallHttpMethod(method string, url string, headers map[string]string, body io.Reader, disableSecurityCheck bool) (res *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println("Error in Call HTTP")
		return
	}
	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
	var client *http.Client
	if disableSecurityCheck {
		client = httpClientDisabledSecurityCheck
	} else {
		client = httpClientEnabledSecurityCheck
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return resp, err
}

func ExecuteCurl(apiUrl string, accessToken string, path string) ([]map[string]interface{}, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["Authorization"] = "bearer " + accessToken
	var decodedBodyArray []map[string]interface{}
	hasNextUrl := true
	url := apiUrl + path
	for hasNextUrl {
		curlResponse, err := CallHttpMethod("GET", url, headers, nil, true)
		defer curlResponse.Body.Close()
		var decodedBody map[string]interface{}
		if err != nil {
			fmt.Printf("Error while CURL call")
			return decodedBodyArray, err
		} else {
			bodyBytes, err2 := ioutil.ReadAll(curlResponse.Body)
			if err2 != nil {
				fmt.Printf("Error while decoding curl response")
				return decodedBodyArray, err2
			}

			err = json.Unmarshal(bodyBytes, &decodedBody)
			if err != nil {
				return decodedBodyArray, err
			}
			nextUrl := decodedBody["next_url"]
			if nextUrl != nil {
				nextUrl := decodedBody["next_url"].(string)
				url = apiUrl + nextUrl
			} else {
				hasNextUrl = false
			}
			decodedBodyArray = append(decodedBodyArray, decodedBody)

		}
	}
	return decodedBodyArray, nil
}

func GetAccessToken(loginUrl string, refreshToken string, grantType string) (string, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Accept"] = "application/json"
	headers["Authorization"] = "Basic Y2Y6"

	data := "grant_type=" + grantType + "&client_id=cf&client_secret=&refresh_token=" + refreshToken

	tokenResponse, err := CallHttpMethod("POST", loginUrl+"/oauth/token", headers, strings.NewReader(data), true)
	defer tokenResponse.Body.Close()
	if err != nil {
		fmt.Printf("Error while getting access-token (CF-Login-CURL call)")
		return "", err
	} else {
		bodyBytes, err2 := ioutil.ReadAll(tokenResponse.Body)
		if err2 != nil {
			fmt.Printf("Error while decoding curl response")
			return "", err2
		}
		var s map[string]string
		err = json.Unmarshal(bodyBytes, &s)
		accessToken := s["access_token"]
		return accessToken, nil
	}
}

func (c *EventCommand) ListEvents(cliConnection plugin.CliConnection, noInstanceNames bool, action string) {
	Initialize()
	fmt.Println("Getting the list of instance events in the org", AddColor(helper.GetOrgName(helper.ReadConfigJsonFile()), constants.Cyan), "/ space", AddColor(helper.GetSpaceName(helper.ReadConfigJsonFile()), constants.Cyan), "...")
	var cmd string
	var actionType string
	var userSpaceGuid = helper.GetSpaceGUID(helper.ReadConfigJsonFile())

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetHeaderLine(false)
	table.SetAutoFormatHeaders(false)

	table.SetHeader([]string{AddColor("instance_name", constants.White), AddColor("instance_guid", constants.White), AddColor("event_type", constants.White), AddColor("user", constants.White), AddColor("created_at", constants.White)})

	var columns int = 5
	var field = make([]string, columns)

	if action == "create" {
		actionType = "audit.service_instance.create"
		cmd = "/v2/events?q=type:" + actionType + "%3Bspace_guid:" + userSpaceGuid
	} else if action == "update" {
		actionType = "audit.service_instance.update"
		cmd = "/v2/events?q=type:" + actionType + "%3Bspace_guid:" + userSpaceGuid
	} else if action == "delete" {
		actionType = "audit.service_instance.delete"
		cmd = "/v2/events?q=type:" + actionType + "%3Bspace_guid:" + userSpaceGuid
	} else {
		cmd = "/v2/events?q=type+IN+audit.service_instance.delete,audit.service_instance.create,audit.service_instance.update%3Bspace_guid:" + userSpaceGuid
	}
	var AuthorizationEndpoint string = helper.GetLoginEndpoint(helper.ReadConfigJsonFile())
	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var refreshToken string = helper.GetRefreshToken(helper.ReadConfigJsonFile())

	accessToken, _ := GetAccessToken(AuthorizationEndpoint, refreshToken, "refresh_token")
	curlResponse, err := ExecuteCurl(apiEndpoint, accessToken, cmd)

	if err != nil {
		fmt.Println(AddColor("FAILED", constants.Red))
		fmt.Printf("Errors in getting service instance events. ", err)
		return
	} else {
		fmt.Println(AddColor("OK", constants.Green))
		for _, val := range curlResponse {
			resources := val["resources"].([]interface{})
			for _, resource := range resources {
				resourceObj := resource.(map[string]interface{})
				resourceObjMetadata := resourceObj["metadata"].(map[string]interface{})
				field[4] = resourceObjMetadata["created_at"].(string)
				resourceObjEntity := resourceObj["entity"].(map[string]interface{})
				field[1] = resourceObjEntity["actee"].(string)
				field[2] = resourceObjEntity["type"].(string)
				field[3] = resourceObjEntity["actor_name"].(string)
				field[0] = resourceObjEntity["actee_name"].(string)
				field[1] = AddColor(field[1], constants.Cyan)
				table.Append(field)
			}
		}
	}
	table.Render()
}
