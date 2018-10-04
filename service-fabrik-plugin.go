package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/cloudfoundry-incubator/service-fabrik-cli-plugin/backup"
	"github.com/cloudfoundry-incubator/service-fabrik-cli-plugin/errors"
	"github.com/cloudfoundry-incubator/service-fabrik-cli-plugin/events"
	"github.com/cloudfoundry-incubator/service-fabrik-cli-plugin/helper"
	"github.com/cloudfoundry-incubator/service-fabrik-cli-plugin/restore"
	"github.com/cloudfoundry/cli/cf/trace"
	"io"
	"os"
	"strconv"
	"strings"
)

//Dynamically set during build time
var Version string = "0.0.0"

type ServiceFabrikPlugin struct {
	cliConnection        plugin.CliConnection
	stdout               io.Writer
	terminalOutputSwitch TerminalOutputSwitch
	logger               trace.Printer
}

type TerminalOutputSwitch interface {
	DisableTerminalOutput(bool)
}

func (cmd *ServiceFabrikPlugin) DisableTerminalOutput(disable bool, retVal *bool) error {
	cmd.terminalOutputSwitch.DisableTerminalOutput(disable)
	*retVal = true
	return nil
}

func main() {
	plugin.Start(new(ServiceFabrikPlugin))
}

func (serviceFabrikPlugin *ServiceFabrikPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	argLength := len(args) // Whatever comes after the "cf" word as command are part of args.
	helper.CreateConfFile()

	//Display help text if user enters "cf backup"
	if argLength == 1 && args[0] == "backup" {
		serviceFabrikPlugin.printHelp(false, "")
		return
	}

	//Display help text if user enters "cf restore"
	if argLength == 1 && args[0] == "restore" {
		serviceFabrikPlugin.printHelp(false, "")
		return
	}

	if args[0] == "backup" { //If user enters, "cf backup [BACKUP ID]"
		if argLength != 2 {
			errors.IncorrectNumberOfArguments()
			return
		}
		backup.NewBackupCommand(cliConnection).BackupInfo(cliConnection, args[1])
	}

	var cmds []string = strings.Split(args[0], "-")
	//3 overall switches: backup, restore & events
	if len(cmds) == 2 {

		switch cmds[1] {
		case "backup":
			if argLength > 3 {
				errors.IncorrectNumberOfArguments()
			} //Error code applicable to all backup commands.
			//Internally split into start, abort, list, delete
			switch cmds[0] {
			case "start":
				if argLength != 2 {
					errors.IncorrectNumberOfArguments()
				}
				fmt.Println("Are you sure you want to start backup? (y/n)")
				var userChoice string
				fmt.Scanln(&userChoice)
				if userChoice == "y" {
					backup.NewBackupCommand(cliConnection).StartBackup(cliConnection, args[1])
				} else {
					os.Exit(7)
				}
			case "abort":
				if argLength != 2 {
					errors.IncorrectNumberOfArguments()
				}
				fmt.Println("Are you sure you want to abort backup? (y/n)")
				var userChoice string
				fmt.Scanln(&userChoice)
				if userChoice == "y" {
					backup.NewBackupCommand(cliConnection).AbortBackup(cliConnection, args[1])
				} else {
					os.Exit(7)
				}

			//List backup has 2 criteria: listing all backups in space and/or listing all backups of the service-instance name given by user.
			case "list":
				if argLength == 2 {
					if args[1] == "--no-name" {
						backup.NewBackupCommand(cliConnection).ListBackups(cliConnection, true)
					} else {
						backup.NewBackupCommand(cliConnection).ListBackupsByInstance(cliConnection, args[1], "", false)
					}
				}
				if argLength == 1 {
					backup.NewBackupCommand(cliConnection).ListBackups(cliConnection, false)
				}
				if argLength == 3 {
					if args[2] == "--deleted" {
						backup.NewBackupCommand(cliConnection).ListBackupsByDeletedInstanceName(cliConnection, args[1])
					} else if args[1] == "--guid" {
						backup.NewBackupCommand(cliConnection).ListBackupsByInstance(cliConnection, "", args[2], true)
					} else {
						errors.InvalidArgument()
						serviceFabrikPlugin.printHelp(false, "")
						return
					}
				}
			case "delete":
				if argLength != 2 {
					errors.IncorrectNumberOfArguments()
				}
				//Retrieve user_space_guid
				fmt.Println("Are you sure you want to delete backup? (y/n)")
				var userChoice string
				fmt.Scanln(&userChoice)
				if userChoice == "y" {
					backup.NewBackupCommand(cliConnection).DeleteBackup(cliConnection, args[1])
				} else {
					os.Exit(7)
				}
			}

		case "restore":
			//Internally split into start and abort.
			switch cmds[0] {
			case "start":
				if argLength == 4 {
					if args[2] == "--backup-guid" {
						fmt.Println("Are you sure you want to start restore? (y/n)")
						var userChoice string
						fmt.Scanln(&userChoice)
						if userChoice == "y" {
							restore.NewRestoreCommand(cliConnection).StartRestore(cliConnection, args[1], args[3], "", true, false, false, "", "", false)
						} else {
							os.Exit(7)
						}
					} else if args[2] == "--timestamp" {
						fmt.Println("Are you sure you want to start restore? (y/n)")
						var userChoice string
						fmt.Scanln(&userChoice)
						if userChoice == "y" {
							restore.NewRestoreCommand(cliConnection).StartRestore(cliConnection, args[1], "", args[3], false, false, false, "", "", false)
						} else {
							os.Exit(7)
						}
					} else {
						serviceFabrikPlugin.printHelp(true, "start-restore")
						errors.InvalidArgument()
					}
				} else if argLength == 6 {
					if args[4] == "--source-id" {
						fmt.Println("Are you sure you want to start restore? (y/n)")
						var userChoice string
						fmt.Scanln(&userChoice)
						if userChoice == "y" {
							restore.NewRestoreCommand(cliConnection).StartRestore(cliConnection, args[1], "", args[3], false, true, false, args[5], "", false)
						} else {
							os.Exit(7)
						}
					} else if args[4] == "--source" {
						fmt.Println("Are you sure you want to start restore? (y/n)")
						var userChoice string
						fmt.Scanln(&userChoice)
						if userChoice == "y" {
							restore.NewRestoreCommand(cliConnection).StartRestore(cliConnection, args[1], "", args[3], false, false, true, "", args[5], false)
						} else {
							os.Exit(7)
						}
					} else {
						serviceFabrikPlugin.printHelp(true, "start-restore")
						errors.InvalidArgument()
					}
				} else if argLength == 7 {
					if args[6] == "--deleted" && args[4] == "--source" {
						fmt.Println("Are you sure you want to start restore? (y/n)")
						var userChoice string
						fmt.Scanln(&userChoice)
						if userChoice == "y" {
							restore.NewRestoreCommand(cliConnection).StartRestore(cliConnection, args[1], "", args[3], false, false, false, "", args[5], true)
						} else {
							os.Exit(7)
						}
					} else {
						serviceFabrikPlugin.printHelp(true, "start-restore")
						errors.InvalidArgument()
					}
				} else {
					serviceFabrikPlugin.printHelp(true, "start-restore")
					errors.IncorrectNumberOfArguments()
				}
			case "abort":
				if argLength != 2 {
					errors.IncorrectNumberOfArguments()
				}
				fmt.Println("Are you sure you want to start backup? (y/n)")
				var userChoice string
				fmt.Scanln(&userChoice)
				if userChoice == "y" {
					restore.NewRestoreCommand(cliConnection).AbortRestore(cliConnection, args[1])
				} else {
					os.Exit(7)
				}
			}
		case "events":
			switch cmds[0] {
			case "instance":
				if argLength > 2 {
					errors.IncorrectNumberOfArguments()
				}
				if argLength == 2 {
					if args[1] == "--delete" {
						events.NewEventsCommand(cliConnection).ListEvents(cliConnection, true, "delete")
					} else if args[1] == "--create" {
						events.NewEventsCommand(cliConnection).ListEvents(cliConnection, true, "create")
					} else if args[1] == "--update" {
						events.NewEventsCommand(cliConnection).ListEvents(cliConnection, true, "update")
					} else {
						errors.InvalidArgument()
						serviceFabrikPlugin.printHelp(false, "")
						return
					}
				}
				if argLength == 1 {
					events.NewEventsCommand(cliConnection).ListEvents(cliConnection, true, "")
				}

			}
		}
	}
}

func (c *ServiceFabrikPlugin) printHelp(printSingleCommandBool bool, command string) {
	metadata := c.GetMetadata()
	if printSingleCommandBool == false {
		for _, command := range metadata.Commands {
			fmt.Println("Name:")
			fmt.Printf("    %-s - %-s\n", command.Name, command.HelpText)
			fmt.Println("Usage:")
			fmt.Printf("    %-s\n", command.UsageDetails.Usage)
			fmt.Println()
		}
	} else {
		for _, command := range metadata.Commands {
			if command.Name == "start-restore" {
				fmt.Println("Name:")
				fmt.Printf("    %-s - %-s\n", command.Name, command.HelpText)
				fmt.Println("Usage:")
				fmt.Printf("    %-s\n", command.UsageDetails.Usage)
				fmt.Println()
			}
		}
	}

}

func (serviceFabrikPlugin *ServiceFabrikPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    "ServiceFabrikPlugin",
		Version: setVersion(Version),
		Commands: []plugin.Command{
			{ // required to be a registered command
				Name:     "start-backup",
				HelpText: "Start backup of a service instance",
				UsageDetails: plugin.Usage{
					Usage: "cf start-backup SERVICE_INSTANCE_NAME",
				},
			},
			{
				Name:     "abort-backup",
				HelpText: "Abort backup of a service instance",
				UsageDetails: plugin.Usage{
					Usage: "cf abort-backup SERVICE_INSTANCE_NAME",
				},
			},
			{
				Name:     "list-backup",
				HelpText: "List backup(s) of a service instance",
				UsageDetails: plugin.Usage{
					Usage: "cf list-backup [SERVICE_INSTANCE_NAME] \n    cf list-backup [SERVICE_INSTANCE_NAME] --deleted \n    cf list-backup  --guid INSTANCE_GUID",
				},
			},
			{
				Name:     "instance-events",
				HelpText: "List events for service instances",
				UsageDetails: plugin.Usage{
					Usage: "cf instance-events [--delete|--create|--update]",
				},
			},
			/*{
				Name:     "delete-backup",
				HelpText: "Delete backup of the given BACKUP_ID",
				UsageDetails: plugin.Usage{
					Usage: "cf delete-backup BACKUP_ID",
				},
			},*/
			{
				Name:     "backup",
				HelpText: "Details of the given BACKUP_ID",
				UsageDetails: plugin.Usage{
					Usage: "cf backup BACKUP_ID",
				},
			},
			{
				Name:     "start-restore",
				HelpText: "Start restore of a service instance",
				UsageDetails: plugin.Usage{
					Usage: "cf start-restore SERVICE_INSTANCE_NAME --backup_guid BACKUP_ID \n    cf start-restore SERVICE_INSTANCE_NAME --timestamp TIME_STAMP\n    cf start-restore SERVICE_INSTANCE_NAME --timestamp TIME_STAMP --source SOURCE_SERVICE_INSTANCE_NAME\n    cf start-restore SERVICE_INSTANCE_NAME --timestamp TIME_STAMP --source-id SOURCE_SERVICE_INSTANCE_GUID\n    cf start-restore SERVICE_INSTANCE_NAME --timestamp TIME_STAMP --source SOURCE_SERVICE_INSTANCE_NAME --deleted",
				},
			},
			{
				Name:     "abort-restore",
				HelpText: "Abort restore of a service instance",
				UsageDetails: plugin.Usage{
					Usage: "cf abort-restore SERVICE_INSTANCE_NAME",
				},
			},
		},
	}
}

func setVersion(version string) plugin.VersionType {
	mmb := strings.Split(version, ".")
	if len(mmb) != 3 {
		panic("invalid version: " + version)
	}
	major, _ := strconv.Atoi(mmb[0])
	minor, _ := strconv.Atoi(mmb[1])
	build, _ := strconv.Atoi(mmb[2])

	return plugin.VersionType{
		Major: major,
		Minor: minor,
		Build: build,
	}
}
