package errors

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/fatih/color"
)

func Condition(cond bool, message string) {
	if !cond {
		color.Red("FAILED")
		panic(errors.New("PLUGIN ERROR: " + message))
	}
}

func ErrorIsNil(err error) {
	if err != nil {
		color.Red("FAILED")
		Condition(false, "error not nil, "+err.Error())
	}
}

func IncorrectNumberOfArguments() {
	color.Red("FAILED")
	fmt.Println("You have entered incorrect number of arguments.")
	fmt.Println("Enter 'cf backup' to check the list of commands and their usage.")
	os.Exit(1)
}

func InstanceGuidNotFound(instanceName string) {
	color.Red("FAILED")
	fmt.Println("Instance Guid not found for the given deleted instance " + instanceName + ".")
	fmt.Println("Enter 'cf backup' to check the list of commands and their usage.")
	os.Exit(1)
}

func InvalidArgument() {
	color.Red("FAILED")
	fmt.Println("You have entered an invalid argument.")
	fmt.Println("Enter 'cf backup' to check the list of commands and their usage.")
	os.Exit(1)
}

func IncorrectSpace(orgName string, spaceName string) {
	color.Red("FAILED")
	fmt.Println("Instance name requested doesn't belong to the org: " + orgName + " and the space: " + spaceName + " Please target the correct org and space.")
	os.Exit(2)
}

func IncorrectInstanceName(instanceName string) {
	color.Red("FAILED")
	fmt.Println("Service Instance \"" + instanceName + "\" doesn't exist.")
	os.Exit(3)
}

func IncorrectServiceType(instanceName string, serviceName string) {
	color.Red("FAILED")
	fmt.Println("Service Instance \"" + instanceName + "\" is of service \"" + serviceName + "\".")
	fmt.Println("Service \"" + serviceName + "\" is not supported for this command.")
	os.Exit(3)
}

func BackupsNotFound(instanceGuid string) {
	color.Red("FAILED")
	fmt.Println("No backups found for the service instance Guid \"" + instanceGuid + "\".")
	os.Exit(3)
}

func CfCliPluginError(temp string) {
	color.Red("FAILED")
	fmt.Println(" PLUGIN ERROR: Error from Cli Command: cf ", temp)
	os.Exit(4)
}

func FileReadingError(filename string) {
	color.Red("FAILED")
	fmt.Println("Encountered error while trying to read the file: ", filename)
	os.Exit(5)
}

func NoAccessTokenError(val string) {
	color.Red("FAILED")
	fmt.Println("No " + val + " was found.")
	fmt.Println("You may be logged out. Please log in to continue.")
	os.Exit(6)
}

func HomeDirNotFound(err error) {
	color.Red("FAILED")
	log.Fatal(err)
	os.Exit(7)
}
