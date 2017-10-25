package guidTranslator

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/SAP/service-fabrik-cli-plugin/errors"
	"github.com/SAP/service-fabrik-cli-plugin/helper"
	"strings"
)

type CliCmd struct{}

func findNextPage(output []string) string {
	for _, val := range output {
		var str []string = strings.Fields(val)

		if len(str) > 1 {
			if strings.Compare(str[0], "\"next_url\":") == 0 {
				next_page_url := str[1]
				next_page_url = strings.TrimRight(next_page_url, ",")
				next_page_url = strings.Trim(next_page_url, "\"")
				return next_page_url
			}
		}
	}
	return "null"
}

func FindInstanceName(cliConnection plugin.CliConnection, InstanceGuid string, output []string) string {
	var cmd string
	var guidTemp string
	cmd = "/v2/service_instances"
	var err error
	var nextPage bool = false

	for cmd != "null" {
		if output == nil || nextPage == true {
			output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
		}
		if err != nil {
			errors.CfCliPluginError(cmd)
		}

		for index, val := range output {
			_ = index
			var str []string = strings.Fields(val)

			if len(str) > 1 {

				//extract guid of service-instance
				if strings.Compare(str[0], "\"guid\":") == 0 {
					guidTemp = str[1]
				}
				//Comaparing the guid with the actual instance guid
				if strings.Contains(guidTemp, InstanceGuid) {
					if strings.Compare(str[0], "\"name\":") == 0 {
						str[1] = strings.TrimRight(str[1], ",")
						return str[1]
					} // returning instance name based on match

				}

			}

		}
		cmd = findNextPage(output)
		if cmd != "null" {
			nextPage = true
		}
	}
	return "" //if no match is found, return "Invalid name"
}

func FindServiceName(cliConnection plugin.CliConnection, serviceId string, output []string) string {
	var cmd string
	var serviceIdTemp string
	var label_temp string
	cmd = "/v2/services"
	var err error
	var nextPage bool = false

	for cmd != "null" {
		if output == nil || nextPage == true {
			output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
		}
		if err != nil {
			errors.CfCliPluginError(cmd)
		}

		for index, val := range output {
			_ = index
			var str []string = strings.Fields(val)

			if len(str) > 1 {

				if strings.Compare(str[0], "\"label\":") == 0 {
					label_temp = str[1]
				}

				//extract guid of service-instance
				if strings.Compare(str[0], "\"unique_id\":") == 0 {
					serviceIdTemp = str[1]
				}
				//Comaparing the guid with the actual instance guid
				if strings.Contains(serviceIdTemp, serviceId) {
					label_temp = strings.TrimRight(label_temp, ",")
					return label_temp

				}

			}

		}
		cmd = findNextPage(output)
		if cmd != "null"{
			nextPage = true
		}
	}

	return "Invalid Name" //if no match is found, return "Invalid name"
}

func FindPlanName(cliConnection plugin.CliConnection, planId string, output []string) string {
	var cmd string
	var planIdTemp string
	var planNameTemp string
	cmd = "/v2/service_plans"
	var err error
	var nextPage bool = false
	
	for cmd != "null" {
		if output == nil || nextPage == true  {
			output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
		}
		if err != nil {
			errors.CfCliPluginError(cmd)
		}

		for index, val := range output {
			_ = index
			var str []string = strings.Fields(val)
			if len(str) > 1 {

				if strings.Compare(str[0], "\"name\":") == 0 {
					planNameTemp = str[1]
				}

				//extract guid of service-instance
				if strings.Compare(str[0], "\"unique_id\":") == 0 {
					planIdTemp = str[1]
				}
				//Comaparing the guid with the actual instance guid
				if strings.Contains(planIdTemp, planId) {
					planNameTemp = strings.TrimRight(planNameTemp, ",")
					return planNameTemp
				}

			}

		}
		cmd = findNextPage(output)
		if cmd != "null" {
			nextPage = true
		}
	}

	return "Invalid Name" //if no match is found, return "Invalid name"
}

func FindServiceId(cliConnection plugin.CliConnection, serviceGuid string, output []string) string {
	var cmd string
	var guidTemp string
	var serviceId string
	var nextPage bool = false	

	cmd = "/v2/services"
	var err error

	for cmd != "null" {
		if output == nil || nextPage == true {
			output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
		}
		if err != nil {
			errors.CfCliPluginError(cmd)
		}

		for index, val := range output {
			_ = index
			var str []string = strings.Fields(val)

			if len(str) > 1 {

				//extract guid of service
				if strings.Compare(str[0], "\"guid\":") == 0 {
					guidTemp = str[1]
				}
				//compare the guid with guid of actual service
				if strings.Contains(guidTemp, serviceGuid) {
					if strings.Compare(str[0], "\"unique_id\":") == 0 {
						serviceId = str[1] //set serviceId of the service based on the match
						serviceId = strings.TrimRight(serviceId, ",")
						return serviceId
					}

				}

			}

		}
		cmd = findNextPage(output)
		if cmd != "null" {
			nextPage = true
		}
	}
	return "Invalid servicePlanId"
}

func FindServicePlanId(cliConnection plugin.CliConnection, servicePlanGuid string, output []string) string {
	var cmd string
	var guidTemp string
	var servicePlanId string
	cmd = "/v2/service_plans"
	var err error
	var nextPage bool = false

	for cmd != "null" {
		if output == nil || nextPage == true {
			output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
		}
		if err != nil {
			errors.CfCliPluginError(cmd)
		}

		for index, val := range output {
			_ = index
			var str []string = strings.Fields(val)

			if len(str) > 1 {

				//extract guid of service-plan
				if strings.Compare(str[0], "\"guid\":") == 0 {
					guidTemp = str[1]
				}
				//compare the servicePlanGuid with servicePlanGuid of actual service_plan
				if strings.Contains(guidTemp, servicePlanGuid) {
					if strings.Compare(str[0], "\"unique_id\":") == 0 {
						servicePlanId = str[1]
						servicePlanId = strings.TrimRight(servicePlanId, ",")
						return servicePlanId
					} //set servicePlanId as per match

				}

			}

		}
		cmd = findNextPage(output)
		if cmd != "null"{
			nextPage = true
		}
	}
	return "Invalid servicePlanGuid"
}

func FindServiceGUId(cliConnection plugin.CliConnection, servicePlanGuid string, output []string) string {
	var cmd string
	var guidTemp string
	var serviceGuid string
	cmd = "/v2/service_plans"
	var err error
	var nextPage = false

	for cmd != "null" {
		if output == nil || nextPage == true {
			output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
		}

		if err != nil {
			errors.CfCliPluginError(cmd)
		}

		for index, val := range output {
			_ = index
			var str []string = strings.Fields(val)

			if len(str) > 1 {

				//extract guid of service-plan
				if strings.Compare(str[0], "\"guid\":") == 0 {
					guidTemp = str[1]
				}
				//compare the servicePlanGuid with servicePlanGuid of actual service_plan
				if strings.Contains(guidTemp, servicePlanGuid) {
					if strings.Compare(str[0], "\"service_guid\":") == 0 {
						serviceGuid = str[1]
						return serviceGuid
					} //set serviceGuid as per match

				}

			}

		}
		cmd = findNextPage(output)
		if cmd != "null"{
			nextPage = true
		}
	}
	return "Invalid_Service_Guid"
}

func FindInstanceGuid(cliConnection plugin.CliConnection, instanceName string, output []string, userSpaceGuid string) string {
	var cmd string
	var err error
	cmd = "/v2/service_instances"

	var flag int = 0
	var guidTemp string
	var servicePlanGuidTemp string = "_"
	var instanceNameTemp string
	var spaceGuid string = "_"
	var guid string
	var nextPage bool = false

	//str3 = "\""+args[3]+"\""
	var userInput string = "\"" + instanceName + "\""

	for cmd != "null" {
		if output == nil || nextPage == true {
			output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
		}

		if err != nil {
			errors.CfCliPluginError(cmd)
		}

		for index, val := range output {
			_ = index
			var str []string = strings.Fields(val)

			if len(str) > 1 {

				//extract guid of service-instance
				if strings.Compare(str[0], "\"guid\":") == 0 {
					guidTemp = str[1]
				}

				//extract service-instance-name
				if strings.Compare(str[0], "\"name\":") == 0 {
					instanceNameTemp = str[1]
				}

				//Compare with userInput
				if strings.Contains(instanceNameTemp, userInput) {
					//if true; extract servicePlanGuid & spaceGuid
					flag = 1
					if strings.Compare(str[0], "\"service_plan_guid\":") == 0 {
						servicePlanGuidTemp = str[1]
					}

					if strings.Compare(str[0], "\"space_guid\":") == 0 {
						spaceGuid = str[1]
					}

					//Retrieve userSpaceGuid
					if userSpaceGuid == "" {
						userSpaceGuid = helper.GetSpaceGUID(helper.ReadConfigJsonFile())
					}

					//Compare spaceGuid with userSpaceGuid

					//if true:
					if (servicePlanGuidTemp != "_") && (spaceGuid != "_") {
						if strings.Contains(spaceGuid, userSpaceGuid) {
							guid = guidTemp
							spaceGuid = userSpaceGuid
							return guid
						}

					}

				}
			}
		}
		cmd = findNextPage(output)
		if cmd != "null"{
			nextPage = true
		}
	}
	if flag == 0 {
		errors.IncorrectInstanceName(instanceName)
	}
	errors.IncorrectSpace(helper.GetOrgName(helper.ReadConfigJsonFile()), helper.GetSpaceName(helper.ReadConfigJsonFile()))
	return "Invalid_Instance_Guid"
}

func FindServicePlanGuid(cliConnection plugin.CliConnection, instanceName string, output []string, userSpaceGuid string) string {

	var cmd string
	cmd = "/v2/service_instances"

	//var guidTemp string
	var servicePlanGuidTemp string = "_"
	var instanceNameTemp string
	var spaceGuid string = "_"
	var servicePlanGuid string

	//str3 = "\""+args[3]+"\""
	var userInput string = "\"" + instanceName + "\""

	var flag int = 0
	var err error
	var nextPage bool = false

	for cmd != "null" {
		if output == nil || nextPage == true {
			output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", cmd)
		}

		if err != nil {
			errors.CfCliPluginError(cmd)
		}

		for index, val := range output {
			_ = index
			var str []string = strings.Fields(val)

			if len(str) > 1 {

				//extract guid of service-instance
				//	if(strings.Compare(str[0],"\"guid\":")==0) {
				//guidTemp = str[1]
				//	}

				//extract service-instance-name
				if strings.Compare(str[0], "\"name\":") == 0 {
					instanceNameTemp = str[1]
				}

				//Compare with userInput
				if strings.Contains(instanceNameTemp, userInput) {
					//if true; extract servicePlanGuid & spaceGuid
					flag = 1
					if strings.Compare(str[0], "\"service_plan_guid\":") == 0 {
						servicePlanGuidTemp = str[1]
					}

					if strings.Compare(str[0], "\"space_guid\":") == 0 {
						spaceGuid = str[1]
					}

					//Retrieve userSpaceGuid
					if userSpaceGuid == "" {
						userSpaceGuid = helper.GetSpaceGUID(helper.ReadConfigJsonFile())
					}

					//Compare spaceGuid with userSpaceGuid

					//if true:
					if (servicePlanGuidTemp != "_") && (spaceGuid != "_") {
						if strings.Contains(spaceGuid, userSpaceGuid) {
							servicePlanGuid = servicePlanGuidTemp
							spaceGuid = userSpaceGuid
							return servicePlanGuid
						}

					}

				}

			}
		}
		cmd = findNextPage(output)
		if cmd != "null"{
			nextPage = true
		}
	}
	if flag == 0 {
		errors.IncorrectInstanceName(instanceName)
	}
	errors.IncorrectSpace(helper.GetOrgName(helper.ReadConfigJsonFile()), helper.GetSpaceName(helper.ReadConfigJsonFile()))
	return "Invalid Service_Plan_Guid"
}
