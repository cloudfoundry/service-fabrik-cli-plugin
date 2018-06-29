# ServiceFabrik CF CLI Plugin

This is a [Cloud Foundry CLI](https://github.com/cloudfoundry/cli) plugin for performing various backup and restore operations on [service-instances](https://docs.cloudfoundry.org/devguide/services/managing-services.html) in Cloud Foundry, such as starting/aborting a backup, listing all backups, removing backups, starting/aborting a restore, etc. 

# Developer's Guide

## Pre-requisites:
- **go binaries:** You need to have `go` in your build environment. Inorder to do that, you need to download and install the golang package. The binaries for download and corresponding installation instructions can be found [here](https://golang.org/dl/).
- **CloudFoundry Command Line Interface (CF CLI)**: You need to have CF CLI installed for the plugin to work since it is a plugin built on CF CLI. The installation instructions for CF CLI can be found [here](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html).

## Building the plugin 
Clone this repo and build it. For this execute following commands on Linux or Mac OS X system
```
$ go get github.com/SAP/service-fabrik-cli-plugin
$ cd $GOPATH/src/github.com/SAP/service-fabrik-cli-plugin
$ go build .
```

The above will clone your repo into default $GOPATH. If you want to setup a different $GOPATH and work on that, then execute following commands on a Linux or Mac OS X system:

```
$ mkdir -p service-fabrik-cli-plugin/src/github.com/SAP/
$ export GOPATH=$(pwd)/service-fabrik-cli-plugin:$GOPATH
$ cd service-fabrik-cli-plugin/src/github.com/SAP/
$ git clone https://github.com/SAP/service-fabrik-cli-plugin.git
$ cd service-fabrik-cli-plugin
$ go build .
```
This will generate a binary executable with the name `service-fabrik-cli-plugin`.

## Installation and Getting Started

- Ensure that CF CLI is installed and working. 
- You should have generated a plugin executable by building the package. Refer to previous section for details.
- Before using the plugin, you need to install it to the CLI.

For Windows
```
cf install-plugin C:\Users\[username]\github.com\SAP\service-fabrik-cli-plugin\servicefabrik_cli_plugin_windows_amd64.exe
```
For Mac
```
cf install-plugin ~/github.com/SAP/service-fabrik-cli-plugin/servicefabrik_cli_plugin_darwin_amd64
```
For Linux
```
cf install-plugin ~/github.com/SAP/service-fabrik-cli-plugin/servicefabrik_cli_plugin_linux_amd64
```
The installation instructions given here imply that the working directory is the home directory. Kindly change it to the proper directory structure as given here, if it is not so.

This CF CLI plugin is only available for ServiceFabrik broker, so it can only be used with CF installations in which this service broker is available.
You can also list all available commands and their usage with `cf backup`. For more information, see [Commands](#commands) and [Further Reading](#further_reading) below.

## Building new release version
You can automatically build new release for all supported platforms by calling the build.sh script with the version of the build.
The version will be automatically included in the plugin, so it will be reported by `cf plugins`.

:rotating_light: Note that the version parameter should follow the semver format (e.g. 1.2.3).
```bash
./build.sh 1.2.3
```
This will produce ` servicefabrik_cli _plugin_linux_amd64`, ` servicefabrik_cli_plugin_darwin_amd64` and ` servicefabrik_cli_plugin_windows_amd64` in the repo's root directory.

# Commands

This plugin adds the following commands:

Command Name | Command Description
--- | ---
`cf backup` | Show the list of all commands and their usage.
`cf backup BACKUP_ID` | Show the information about a particular backup.
` cf list-backup ` | Show the list of all backups present in the space.
` cf list-backup SERVICE_INSTANCE_NAME ` | Show the list of all backups for the given service-fabrik service instance.
` cf list-backup --guid SERVICE_INSTANCE_GUID` | Show the list of all backups for the given service-fabrik service instance. The argument has to be the guid of the service instance. (Works even for a deleted instance.)
`cf list-backup SERVICE_INSTANCE_NAME --deleted` | Shows the list of all backups for a deleted service-fabrik service instance. (Works only for a deleted service-instance.)
`cf instance-events` | Lists all events including create, update and delete events triggered for all service instances present in the space.
`cf instance-events --create` | List all create service instance events in the space.
`cf instance-events --update` | List all update service instance events in the space.
`cf instance-events --delete` | List all delete service instance events in the space.
` cf start-restore SERVICE_INSTANCE_NAME BACKUP_ID ` | Start restore of a service-fabrik service instance from the given backup id.
` cf abort-restore SERVICE_INSTANCE_NAME` | Abort restore of a service-fabrik service instance.
 
For more information, see the command help output available via `cf [command] --help` or `cf help [command]`.

# Further Reading
User Documentation: [user_documentation_cf_cli_plugin.md](https://github.com/SAP/service-fabrik-cli-plugin/blob/master/user_documentation_cf_cli_plugin.md)

## How to obtain support

If you need any support, have any question or have found a bug, please report it in the [GitHub bug tracking system](https://github.com/SAP/service-fabrik-cli-plugin/issues). We shall get back to you.

## LICENSE

This project is licensed under the Apache Software License, v. 2 except as noted otherwise in the [LICENSE](LICENSE) file.



