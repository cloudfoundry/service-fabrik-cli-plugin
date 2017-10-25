package restore

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/fatih/color"
	"github.com/SAP/service-fabrik-cli-plugin/errors"
	"github.com/SAP/service-fabrik-cli-plugin/guidTranslator"
	"github.com/SAP/service-fabrik-cli-plugin/helper"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type RestoreCommand struct {
	cliConnection plugin.CliConnection
}

func NewRestoreCommand(cliConnection plugin.CliConnection) *RestoreCommand {
	command := new(RestoreCommand)
	command.cliConnection = cliConnection
	return command
}

const (
	red   color.Attribute = color.FgRed
	green color.Attribute = color.FgGreen
	cyan  color.Attribute = color.FgCyan
	white color.Attribute = color.FgWhite
)

func AddColor(text string, textColor color.Attribute) string {
	printer := color.New(textColor).Add(color.Bold).SprintFunc()
	return printer(text)
}

type Configuration struct {
	ServiceBroker       string
	ServiceBrokerExtUrl string
	SkipSslFlag         bool
}

func GetBrokerName() string {
	return getConfiguration().ServiceBroker
}

func GetExtUrl() string {
	return getConfiguration().ServiceBrokerExtUrl
}

func GetskipSslFlag() bool {
	return getConfiguration().SkipSslFlag
}

func getConfiguration() Configuration {
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

func GetHttpClient() *http.Client {
	//Skip ssl verification.
	client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: GetskipSslFlag()},
                        Proxy:           http.ProxyFromEnvironment,
                },
                Timeout: time.Duration(180) * time.Second,
        }
        return client
}

func GetResponse(client *http.Client, req *http.Request) *http.Response {
	req.Header.Set("Authorization", helper.GetAccessToken(helper.ReadConfigJsonFile()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	errors.ErrorIsNil(err)
	return resp
}

func (c *RestoreCommand) StartRestore(cliConnection plugin.CliConnection, serviceInstanceName string, backupId string) {
	fmt.Println("Starting restore for ", AddColor(serviceInstanceName, cyan), "...")

	if helper.GetAccessToken(helper.ReadConfigJsonFile()) == "" {
		errors.NoAccessTokenError("Access Token")
	}

	client := GetHttpClient()

	var jsonprep string = `{"backup_guid": "` + backupId + `"}`
	var jsonStr = []byte(jsonprep)
	var req_body = bytes.NewBuffer(jsonStr)

	var guid string = guidTranslator.FindInstanceGuid(cliConnection, serviceInstanceName, nil, "")
	guid = strings.TrimRight(guid, ",")
	guid = strings.Trim(guid, "\"")

	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var broker string = GetBrokerName()
	var extUrl string = GetExtUrl()

	apiEndpoint = strings.Replace(apiEndpoint, "api", broker, 1)

	var url string = apiEndpoint + extUrl + "/service_instances/" + guid + "/restore"
	req, err := http.NewRequest("POST", url, req_body)
	var resp *http.Response = GetResponse(client, req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.Status != "202 Accepted" {
		fmt.Println(AddColor("FAILED", red))
		var message string = string(body)
		var parts []string = strings.Split(message, ":")
		fmt.Println(parts[2])
	}

	if resp.Status == "202 Accepted" {
		fmt.Println(AddColor("OK", green))
		fmt.Println("Restore has been initiated for the instance name:", AddColor(serviceInstanceName, cyan), " and from the backup id:", AddColor(backupId, cyan))
		fmt.Println("Please check the status of restore by entering 'cf service SERVICE_INSTANCE_NAME'")
	}

	errors.ErrorIsNil(err)

}

func (c *RestoreCommand) AbortRestore(cliConnection plugin.CliConnection, serviceInstanceName string) {
	fmt.Println("Aborting restore for ", AddColor(serviceInstanceName, cyan), "...")

	if helper.GetAccessToken(helper.ReadConfigJsonFile()) == "" {
		errors.NoAccessTokenError("Access Token")
	}

	client := GetHttpClient()

	var guid string = guidTranslator.FindInstanceGuid(cliConnection, serviceInstanceName, nil, "")
	guid = strings.TrimRight(guid, ",")
	guid = strings.Trim(guid, "\"")

	var userSpaceGuid string = helper.GetSpaceGUID(helper.ReadConfigJsonFile())

	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var broker string = GetBrokerName()
	var extUrl string = GetExtUrl()

	apiEndpoint = strings.Replace(apiEndpoint, "api", broker, 1)

	var url string = apiEndpoint + extUrl + "/service_instances/" + guid + "/restore?space_guid=" + userSpaceGuid
	req, err := http.NewRequest("DELETE", url, nil)

	var resp *http.Response = GetResponse(client, req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if (resp.Status != "202 Accepted") && (resp.Status != "200 OK") {
		fmt.Println(AddColor("FAILED", red))
		var message string = string(body)
		var parts []string = strings.Split(message, ":")
		fmt.Println(parts[2])
	}

	if resp.Status == "202 Accepted" {
		fmt.Println(AddColor("OK", green))
		fmt.Println("Restore has been aborted for the instance name:", color.CyanString(serviceInstanceName))
	}

	if resp.Status == "200 OK" {
		fmt.Println("currently no restore in progress for this service instance")
	}

	errors.ErrorIsNil(err)
}
