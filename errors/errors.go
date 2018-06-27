package errors

import (
	"fmt"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/fatih/color"
	"log"
	"os"
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

func IncorrectGuid(instanceGuid string) {
	color.Red("FAILED")
	fmt.Println("Service Instance Guid \"" + instanceGuid + "\" doesn't exist.")
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
