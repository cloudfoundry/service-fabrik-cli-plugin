package backup

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/SAP/service-fabrik-cli-plugin/errors"
	"github.com/SAP/service-fabrik-cli-plugin/guidTranslator"
	"github.com/SAP/service-fabrik-cli-plugin/helper"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type BackupCommand struct {
	cliConnection plugin.CliConnection
}

func NewBackupCommand(cliConnection plugin.CliConnection) *BackupCommand {
	command := new(BackupCommand)
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

func (c *BackupCommand) BackupInfo(cliConnection plugin.CliConnection, backupId string) {

	fmt.Println("Retrieving information about backup id: ", AddColor(backupId, cyan), "...")

	if helper.GetAccessToken(helper.ReadConfigJsonFile()) == "" {
		errors.NoAccessTokenError("Access Token")
	}

	client := GetHttpClient()

	var userSpaceGuid string = helper.GetSpaceGUID(helper.ReadConfigJsonFile())

	//TODO: This is a workaround to get refreshed jwt token if it is expired, we need to see if this is correct way??
	var cmd string = "/v2/service_instances"
	output, err := cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
	if output != nil {
		output = nil
	}
	if err != nil {
		errors.CfCliPluginError(cmd)
	}

	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var broker string = GetBrokerName()
	var extUrl string = GetExtUrl()

	apiEndpoint = strings.Replace(apiEndpoint, "api", broker, 1)

	var url string = apiEndpoint + extUrl + "/backups/" + backupId + "?space_guid=" + userSpaceGuid

	req, err := http.NewRequest("GET", url, nil)

	var resp *http.Response = GetResponse(client, req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetHeaderLine(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{" ", " "})

	if resp.Status == "200 OK" {
		fmt.Println(AddColor("OK", green))

		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Invalid response for the request ", err)
		}

		field := response["service_id"].(string)
		table.Append([]string{"service-name", strings.Trim(guidTranslator.FindServiceName(cliConnection, field, nil), "\"")})

		field = response["plan_id"].(string)
		table.Append([]string{"plan-name", strings.Trim(guidTranslator.FindPlanName(cliConnection, field, nil), "\"")})

		field = response["instance_guid"].(string)
		table.Append([]string{"instance-name", strings.Trim(guidTranslator.FindInstanceName(cliConnection, field, nil), "\"")})

		table.Append([]string{"organization-name", helper.GetOrgName(helper.ReadConfigJsonFile())})
		table.Append([]string{"space-name", helper.GetSpaceName(helper.ReadConfigJsonFile())})
		table.Append([]string{"username", response["username"].(string)})
		table.Append([]string{"operation", response["operation"].(string)})
		table.Append([]string{"type", response["type"].(string)})
		table.Append([]string{"backup_guid", response["backup_guid"].(string)})
		table.Append([]string{"trigger", response["trigger"].(string)})
		table.Append([]string{"state", response["state"].(string)})
		table.Append([]string{"started_at", response["started_at"].(string)})
		if _, flag := response["finished_at"].(string); flag {
			table.Append([]string{"finished_at", response["finished_at"].(string)})
		} else {
			table.Append([]string{"finished_at", "null"})
		}
		table.Render()
	}

	if resp.Status != "200 OK" {
		fmt.Println(AddColor("FAILED", red))
		var message string = string(body)
		var parts []string = strings.Split(message, ":")
		fmt.Println(parts[2])
	}

}


func (c *BackupCommand) ListBackupsByDeletedInstanceName(cliConnection plugin.CliConnection, serviceInstanceName string) {
	fmt.Println("Getting the list of  backups in the org", AddColor(helper.GetOrgName(helper.ReadConfigJsonFile()), cyan), "/ space", AddColor(helper.GetSpaceName(helper.ReadConfigJsonFile()), cyan), "/ service instance", AddColor(serviceInstanceName, cyan), "...")

	if helper.GetAccessToken(helper.ReadConfigJsonFile()) == "" {
		errors.NoAccessTokenError("Access Token")
	}

	client := GetHttpClient()

	var userSpaceGuid string = helper.GetSpaceGUID(helper.ReadConfigJsonFile())
	var guid string
	var guidMap map[string]string = guidTranslator.FindDeletedInstanceGuid(cliConnection, serviceInstanceName, nil, "")
	if len(guidMap) > 1 {
		fmt.Println(AddColor("FAILED", red))
		fmt.Println("" + serviceInstanceName + " maps to multiple instance GUIDs, please use 'cf instance-events --delete' to list all instance delete events, get instance guid from list and then use cf list-backup --guid GUID to get details")
		fmt.Println("Enter 'cf backup' to check the list of commands and their usage.")
		os.Exit(1)
	} else {
		for k, _ := range guidMap {
			guid = k
			guid = strings.Trim(guid, ",")
			guid = strings.Trim(guid, "\"")
		}
	}
	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var broker string = GetBrokerName()
	var extUrl string = GetExtUrl()

	apiEndpoint = strings.Replace(apiEndpoint, "api", broker, 1)

	req, err := http.NewRequest("GET", apiEndpoint+extUrl+"/backups"+"?space_guid="+userSpaceGuid+"&instance_id="+guid, nil)
	errors.ErrorIsNil(err)
	var resp *http.Response = GetResponse(client, req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetHeaderLine(false)
	table.SetAutoFormatHeaders(false)

	if resp.Status == "200 OK" {
		fmt.Println(AddColor("OK", green))
		var flag bool

		table.SetHeader([]string{AddColor("backup_guid", white), AddColor("username", white), AddColor("type", white), AddColor("trigger", white), AddColor("started_at", white), AddColor("finished_at", white)})

		var response []interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Invalid response for the request ", err)
		}

		var no_of_columns int = 6
		var field = make([]string, no_of_columns)

		for backup := range response {
			if strings.Contains(guid, (response[backup].(map[string]interface{}))["instance_guid"].(string)) {
				field[1] = (response[backup].(map[string]interface{}))["username"].(string)
				field[2] = (response[backup].(map[string]interface{}))["type"].(string)
				field[0] = (response[backup].(map[string]interface{}))["backup_guid"].(string)
				field[0] = AddColor(field[0], cyan)
				field[3] = (response[backup].(map[string]interface{}))["trigger"].(string)
				field[4] = (response[backup].(map[string]interface{}))["started_at"].(string)
				field[5], flag = (response[backup].(map[string]interface{}))["finished_at"].(string)
				if flag == false {
					field[5] = "null"
				}
				table.Append(field)
			}
		}

	}
	table.Render()
	if resp.Status != "200 OK" {
		fmt.Println(AddColor("FAILED", red))
		var message string = string(body)
		var parts []string = strings.Split(message, ":")
		fmt.Println(parts[2])
	}

}

func (c *BackupCommand) ListBackupsByInstance(cliConnection plugin.CliConnection, serviceInstanceName string, instanceGuid string, guidBool bool) {

	if helper.GetAccessToken(helper.ReadConfigJsonFile()) == "" {
		errors.NoAccessTokenError("Access Token")
	}

	client := GetHttpClient()

	var userSpaceGuid string = helper.GetSpaceGUID(helper.ReadConfigJsonFile())
	var guid string
	if guidBool == false {
		guid = guidTranslator.FindInstanceGuid(cliConnection, serviceInstanceName, nil, "")
		guid = strings.Trim(guid, ",")
		guid = strings.Trim(guid, "\"")
	 fmt.Println("Getting the list of  backups in the org", AddColor(helper.GetOrgName(helper.ReadConfigJsonFile()), cyan), "/ space", AddColor(helper.GetSpaceName(helper.ReadConfigJsonFile()), cyan), "/ service instance", AddColor(serviceInstanceName, cyan), "...")
	} else {
		guid = instanceGuid
	fmt.Println("Getting the list of  backups in the org", AddColor(helper.GetOrgName(helper.ReadConfigJsonFile()), cyan), "/ space", AddColor(helper.GetSpaceName(helper.ReadConfigJsonFile()), cyan), "/ service instance GUID", AddColor(instanceGuid, cyan), "...")
	}
	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var broker string = GetBrokerName()
	var extUrl string = GetExtUrl()

	apiEndpoint = strings.Replace(apiEndpoint, "api", broker, 1)

	req, err := http.NewRequest("GET", apiEndpoint+extUrl+"/backups"+"?space_guid="+userSpaceGuid+"&instance_id="+guid, nil)
	errors.ErrorIsNil(err)

	var resp *http.Response = GetResponse(client, req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetHeaderLine(false)
	table.SetAutoFormatHeaders(false)
	var response []interface{}
	
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Invalid response for the request ", err)
	}

	if (len(response) == 0) && (guidBool == true) {
		errors.IncorrectGuid(guid)
	}

	if resp.Status == "200 OK" {
		fmt.Println(AddColor("OK", green))

		var flag bool

		table.SetHeader([]string{AddColor("backup_guid", white), AddColor("username", white), AddColor("type", white), AddColor("trigger", white), AddColor("started_at", white), AddColor("finished_at", white)})

		var no_of_columns int = 6
		var field = make([]string, no_of_columns)
		for backup := range response {
			field[1] = (response[backup].(map[string]interface{}))["username"].(string)
			field[2] = (response[backup].(map[string]interface{}))["type"].(string)
			field[0] = (response[backup].(map[string]interface{}))["backup_guid"].(string)
			field[0] = AddColor(field[0], cyan)
			field[3] = (response[backup].(map[string]interface{}))["trigger"].(string)
			field[4] = (response[backup].(map[string]interface{}))["started_at"].(string)
			field[5], flag = (response[backup].(map[string]interface{}))["finished_at"].(string)
			if flag == false {
				field[5] = "null"
			}
			table.Append(field)
		}

	} else {
		fmt.Println(AddColor("FAILED", red))
		fmt.Println("Error is here")
		var message string = string(body)
		var parts []string = strings.Split(message, ":")
		fmt.Println(parts[2])
	}
	table.Render()
}

func (c *BackupCommand) ListBackups(cliConnection plugin.CliConnection, noInstanceNames bool) {
	fmt.Println("Getting the list of  backups in the org", AddColor(helper.GetOrgName(helper.ReadConfigJsonFile()), cyan), "/ space", AddColor(helper.GetSpaceName(helper.ReadConfigJsonFile()), cyan), "...")

	if helper.GetAccessToken(helper.ReadConfigJsonFile()) == "" {
		errors.NoAccessTokenError("Access Token")
	}

	client := GetHttpClient()

	var cmd string = "/v2/service_instances"

	output, err := cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
	if output != nil {
		output = nil
	}

	if err != nil {
		errors.CfCliPluginError(cmd)
	}

	var userSpaceGuid string = helper.GetSpaceGUID(helper.ReadConfigJsonFile())

	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var broker string = GetBrokerName()
	var extUrl string = GetExtUrl()

	apiEndpoint = strings.Replace(apiEndpoint, "api", broker, 1)

	req, err := http.NewRequest("GET", apiEndpoint+extUrl+"/backups"+"?space_guid="+userSpaceGuid, nil)
	errors.ErrorIsNil(err)

	var resp *http.Response = GetResponse(client, req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetHeaderLine(false)
	table.SetAutoFormatHeaders(false)
	table.SetColWidth(40)

	if resp.Status == "200 OK" {
		fmt.Println(AddColor("OK", green))

		var response []interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Invalid response for the request ", err)
		}

		var no_of_columns int = 8
		var field = make([]string, no_of_columns)

		if noInstanceNames == true {
			table.SetHeader([]string{AddColor("backup_guid", white), AddColor("instance_guid", white), AddColor("username", white), AddColor("type", white), AddColor("trigger", white), AddColor("started_at", white), AddColor("finished_at", white), AddColor(" ", white)})
		} else {
			table.SetHeader([]string{AddColor("backup_guid", white), AddColor("instance_name", white), AddColor("username", white), AddColor("type", white), AddColor("trigger", white), AddColor("started_at", white), AddColor("finished_at", white), AddColor(" ", white)})
		}

		for backup := range response {
			instance_guid := (response[backup].(map[string]interface{}))["instance_guid"].(string)

			if noInstanceNames == true {
				field[1] = instance_guid
			} else {
				var InstanceName = guidTranslator.FindInstanceName(cliConnection, instance_guid, nil)
				field[1] = strings.Trim(InstanceName, "\"")
				if InstanceName == "" {
					field[7] = "Status: Instance already deleted"
				} else {
					field[7] = ""
				}
			}
			field[2] = (response[backup].(map[string]interface{}))["username"].(string)
			field[3] = (response[backup].(map[string]interface{}))["type"].(string)
			field[0] = (response[backup].(map[string]interface{}))["backup_guid"].(string)
			field[0] = AddColor(field[0], cyan)
			field[4] = (response[backup].(map[string]interface{}))["trigger"].(string)
			field[5] = (response[backup].(map[string]interface{}))["started_at"].(string)
			_, flag := (response[backup].(map[string]interface{}))["finished_at"].(string)
			if flag == false {
				field[6] = "null"
			} else {
				field[6] = (response[backup].(map[string]interface{}))["finished_at"].(string)
			}
			table.Append(field)
		}
	}

	table.Render()
	if resp.Status != "200 OK" {
		fmt.Println(AddColor("FAILED", red))
		var message string = string(body)
		var parts []string = strings.Split(message, ":")
		fmt.Println(parts[2])
	}

}

func (c *BackupCommand) DeleteBackup(cliConnection plugin.CliConnection, backupId string) {
	fmt.Println("Deleting backup for ", AddColor(backupId, cyan), "...")

	if helper.GetAccessToken(helper.ReadConfigJsonFile()) == "" {
		errors.NoAccessTokenError("Access Token")
	}

	client := GetHttpClient()

	var userSpaceGuid string = helper.GetSpaceGUID(helper.ReadConfigJsonFile())

	var cmd string = "/v2/service_instances"

	output, err := cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
	if output != nil {
		output = nil
	}

	if err != nil {
		errors.CfCliPluginError(cmd)
	}

	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var broker string = GetBrokerName()
	var extUrl string = GetExtUrl()

	apiEndpoint = strings.Replace(apiEndpoint, "api", broker, 1)

	var url string = apiEndpoint + extUrl + "/backups/" + backupId + "?space_guid=" + userSpaceGuid
	req, err := http.NewRequest("DELETE", url, nil)

	var resp *http.Response = GetResponse(client, req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp.Status != "200 OK" {
		AddColor("FAILED", red)
		var message string = string(body)
		var parts []string = strings.Split(message, ":")
		fmt.Println(parts[2])
	}
	if resp.Status == "200 OK" {
		AddColor("OK", green)
		fmt.Println("The corresponding backup dataset has been deleted.")
	}
}

func (c *BackupCommand) AbortBackup(cliConnection plugin.CliConnection, serviceInstanceName string) {
	fmt.Println("Aborting backup for ", AddColor(serviceInstanceName, cyan), "...")

	if helper.GetAccessToken(helper.ReadConfigJsonFile()) == "" {
		errors.NoAccessTokenError("Access Token")
	}

	client := GetHttpClient()

	var guid string = guidTranslator.FindInstanceGuid(cliConnection, serviceInstanceName, nil, "")
	guid = strings.TrimRight(guid, ",")
	guid = strings.Trim(guid, "\"")

	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var broker string = GetBrokerName()
	var extUrl string = GetExtUrl()

	apiEndpoint = strings.Replace(apiEndpoint, "api", broker, 1)

	var url string = apiEndpoint + extUrl + "/service_instances/" + guid + "/backup"
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
		fmt.Println("Check the state of the backup using cf backup BACKUP_ID command.")
	}

	if resp.Status == "200 OK" {
		fmt.Println("currently no backup in progress for this service instance")
	}

	errors.ErrorIsNil(err)
}

func (c *BackupCommand) StartBackup(cliConnection plugin.CliConnection, serviceInstanceName string) {
	fmt.Println("Triggering backup for ", AddColor(serviceInstanceName, cyan), "...")

	if helper.GetAccessToken(helper.ReadConfigJsonFile()) == "" {
		errors.NoAccessTokenError("Access Token")
	}

	client := GetHttpClient()

	var jsonprep string = `{"type": "online"}`

	var jsonStr = []byte(jsonprep)
	var req_body = bytes.NewBuffer(jsonStr)

	var guid string = guidTranslator.FindInstanceGuid(cliConnection, serviceInstanceName, nil, "")
	guid = strings.TrimRight(guid, ",")
	guid = strings.Trim(guid, "\"")

	var apiEndpoint string = helper.GetApiEndpoint(helper.ReadConfigJsonFile())
	var broker string = GetBrokerName()
	var extUrl string = GetExtUrl()

	apiEndpoint = strings.Replace(apiEndpoint, "api", broker, 1)

	var url string = apiEndpoint + extUrl + "/service_instances/" + guid + "/backup"

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
		var response = string(body)
		var response_str []string = strings.Split(response, ":")
		response_str[2] = strings.TrimRight(response_str[2], "}")
		fmt.Println("BACKUP_ID is", AddColor(response_str[2], cyan))
		fmt.Println("Check the state of the backup using cf backup BACKUP_ID command.")
	}

	errors.ErrorIsNil(err)

}
