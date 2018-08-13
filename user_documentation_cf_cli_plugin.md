# User Documentation for ServiceFabrik CF CLI Plugin



## Version:

1.1

## Authors:

Pritish Mishra

Subhankar Chattopadhyay

Mayank Tiwary

## Table of Contents:

1. [Introduction](#introduction)
1. [Building/Installing the plugin](#building-and-installing-the-plugin)
1. [Important parameters](#important-parameters)
   1. [SERVICE\_INSTANCE\_NAME](#important-parameters)
   1. [BACKUP\_ID](#important-parameters)
1. [Commands and their usage](#commands-and-their-usage)
   1. [Listing all backups](#listing-all-backups)
   1. [Listing all backups of a service-instance](#listing-all-backups-of-a-service-instance)
   1. [Listing all backups of a deleted service-instance](#listing-deleted-backups)
   1. [Listing all backups of a service-instance by instance guid](#listing-deleted-backups-guid)
   1. [Listing service instance events](#listing-instance-events)
   1. [Starting a restore](#starting-a-restore)
   1. [Aborting a restore](#aborting-a-restore)
1. [Error status](#error-status)
   1. [Unauthorized](#unauthorized)
   1. [Another concurrent operation](#another-concurrent-operation)
   1. [No restore in progress](#no-restore-in-progress)
   1. [Unauthorized](#unauthorized)
1. [Common Failure Scenarios]((#common-failure-scenarios))
   1. [IncorrectNumberOfArguments](#incorrect-number-of-arguments)
   1. [IncorrectSpace](#incorrect-space)
   1. [IncorrectInstanceName](#incorrect-instance-name)
   1. [IncorrectCommandUsage](#incorrect-command-usage)
   1. [UserLoggedOutError](#user-logged-out-error)
   1. [MultipleGUIDError](#multtple-guid-error)



## [Introduction](#introduction)

Service Fabrik CF Cli plugin is a Cloud Foundry CLI plugin for performing various backup and restore operations on service-instances in Cloud Foundry, such as starting/aborting a backup, listing all backups, removing backups, starting/aborting a restore, etc.

This CF CLI plugin is only available for ServiceFabrik broker, so it can only be used with CF installations in which this service broker is available. You can also list all available commands and their usage with &#39;_cf__backup_&#39;.

## [Building and Installing the plugin](#building-and-installing-the-plugin)

The following steps need to be followed for the purpose of building the plugin.

1. Clone this repo and build it. For this execute following commands on Linux or Mac OS X system.
  1. `$ go get github.com/cloudfoundry-incubator/service-fabrik-cli-plugin`
  2. `$ cd $GOPATH/src/github.com/cloudfoundry-incubator/service-fabrik-cli-plugin`
  3. `$ go build .`
2. The above will clone your repo into default $GOPATH. If you want to setup a different $GOPATH and work on that, then execute following commands on a Linux or Mac OS X system:
  1. `$ mkdir -p service-fabrik-cli-plugin/src/github.com/cloudfoundry-incubator/`
  2. `$ export GOPATH=$(pwd)/service-fabrik-cli-plugin:$GOPATH`
  3. `$ cd service-fabrik-cli-plugin/src/github.com/cloudfoundry-incubator/`
  4. `$ git clone https://github.com/cloudfoundry-incubator/service-fabrik-cli-plugin.git`
  5. `$ cd service-fabrik-cli-plugin`
  6. `$ go build .`


The following steps need to be followed for the purpose of installing the plugin.

1. You must be having a built plugin executable by following steps mentioned in previous section.
2. Open the command prompt/terminal.
3. For the purposes of re-installing or fresh install of the plugin:
  
    `cf install-plugin [Built plugin location]/service-fabrik-cli-plugin` [Linux/Mac]
  
    `cf install-plugin [Built plugin location]\service-fabrik-cli-plugin.exe` [Windows]
  
4. Check if the plugin has been successfully installed: cf plugins
  If the installation was successful, you will find &quot;ServiceFabrikPlugin&quot; present in the list of plugins along with the commands it supports.
5. If you have a previous version of the plugin already installed on your system, you need to uninstall it. You can do so by:
  `cf uninstall-plugin ServiceFabrikPlugin`

Then repeat the steps 3 &amp; 4.

## [Important parameters](#important-parameters)

There are two very important parameters which the plugin expects from the user for some or the other command. You might have come across these parameters during the next sections. Sometimes, it may become difficult for the user to understand what these parameters are and hence, may lead to unnecessary errors/inconvenience. Hence, we provide a brief description of these parameters.

### SERVICE\_INSTANCE\_NAME:

This is the name of the service instance created by you. If there is a service, e.g. Blueprint service, you may create an instance of it by specifying the service name and choosing a plan of the service. If you have created such a service instance or are using such an instance, and want to take backup/restore of this instance, you need to provide its name wherever SERVICE\_INSTANCE\_NAME is required by the plugin.

The name should be provided without any quotes. The parameter is case-sensitive and hence, must be exactly the same as the name. If you are not sure what the name of the service-instance is or what service-instances you have created/are using, enter &quot;cf s&quot; on your cli to see a list of service-instances.

### BACKUP\_ID:

This is the unique id of the backup created for a service-instance. As we don&#39;t provide the functionality of providing a name to a backup, each backup can only be identified by its id. Hence, you need to take a note of the backup-id when you create a backup of the instance. You can also check the list of all the backups with the ids for all service-instances or for a specific instance. You can then note this backup-id to use for various commands like deleting a backup or creating a restore from the backup.

The id should be provided without any quotes. The parameter is case-sensitive and hence, must be exactly the same as the id.

## [Commands and their usage](#commands-and-their-usage)

The plugin primarily supports 2 operations: Backup &amp; Restore. This means you can take backup of a service instance and can restore a service instance from this backed-up state. All other functionalities have been designed to facilitate these two operations, such as, Listing all the backups you have taken so far, Deleting a backup, etc. In this section, we discuss all the commands supported by the plugin, their usage and the expected output for a successful execution.


### Listing all backups:

**Command:** cf list-backup

**Usage:** This command is used to list all the backups taken by you, of **any** SERVICE\_INSTANCE. You needn&#39;tprovide any additional parameters. Upon successful execution, the plugin will display a list of all backups. The list will also contain the name and id of service-instance for which each backup was taken.

**Expected Output:**

Getting the list of backups for the org: [ORG NAME] and space: [SPACE NAME]

OK

[List of backups]

**Additional note:** When you do a cf-login, you mention the api-endpoint, org-name and space-name. This plugin displays the list of backups taken for the service-instances present in the space, you have logged in to. Hence, the message says &quot;the backups of the targeted space&quot;. The list of backups, sometimes, can be lengthy. If you want to know the list of backups for a particular service-instance, please refer to the next section.

### Listing all backups of a service-instance:

**Command:** cf list-backup SERVICE\_INSTANCE\_NAME

**Usage:** This command is used to list all the backups taken by you, of a **specific** service-instance. You need to provide the name of service-instance, you want the list of backups for, as a parameter. Upon successful execution, the plugin will display a list of all backups of the service-instance.

**Expected Output:**

Getting the list of backups for the org: [ORG NAME] and space: [SPACE NAME]

OK

[List of backups]

**Example:**

<img src="https://github.com/SAP/service-fabrik-cli-plugin/blob/master/images/plugin_screenshot1.png" height="400">

**Additional note:** This command works exactly the same wayas &quot;listing all backups&quot;. It just filters the listto display only the backups of the service-instance required by you.

### Listing all backups of a deleted service-instance

**Command:** cf list-backup SERVICE\_INSTANCE\_NAME --deleted

**Usage:** This command is used to fetch all backups of a deleted service instance. You need to provide the name of service-instance, you want the list of backups for, as a parameter. Upon successful execution, the plugin will display a list of all backups of the service-instance.
**Expected Output:**

Getting the list of  backups in the org [ORG_NAME] / space [SPACE_NAME] / service instance [SERVICE\_INSTANCE\_NAME] ...

OK

[List of backups]

**Additional note:** This command works same as cf list-backup [SERVICE\_INSTANCE\_NAME], but only works on deleted service instance name. 

### Listing all backups of a service-instance by instance guid:

**Command:** cf list-backup --guid SERVICE\_INSTANCE\_GUID

**Usage:** This command is used to fetch backups for deleted service instances using instance guid. You need to provide the Guid of service-instance, you want the list of backups for, as a parameter. Upon successful execution, the plugin will display a list of all backups of the service-instance.

**Expected Output:**

Getting the list of  backups in the org [ORG_NAME] / space [SPACE_NAME] / service instance GUID [SERVICE\_INSTANCE\_GUID] ...

OK

[List of backups]

**Additional note:** This command works same as cf list-backup [SERVICE\_INSTANCE\_NAME], but also works on deleted service instance. 

### Listing service instance events:

**Command:** cf instance-events [--delete|--create|--update]

**Usage:** This command is used to fetch all events of all service instances within the space. Upon successful execution of the command, the plugin will print all the recorded events of all service instances.

**Expected Output:**

Getting the list of instance events in the org [ORG_NAME] / space [SPACE_NAME] ...

OK

[List of events]

**Additional note:** The successful execution of this command will return all the events releated to all service instances. You can also use flags [--delete|--update|--create] to filter out results based on event type. 

### Starting a restore:

**Command:** cf start-restore SERVICE\_INSTANCE\_NAME BACKUP\_ID

**Usage:** You can restore the state of a service-instance from any previously taken backup. In order to dothat, you must provide the name of the service-instance and the id of the backup from which restore must be done. Upon successful execution, the plugin will initiate the restore process from the backup.

**Expected Output:**

Starting restore

OK

Restore has been initiated for the instance name: [SERVICE\_INSTANCE\_NAME] and from the backup id:

[BACKUP\_ID]

Please check the status of restore by entering &#39;cf service SERVIC\_INSTANCE\_NAME&#39;

**Example:**

<img src="https://github.com/SAP/service-fabrik-cli-plugin/blob/master/images/plugin_screenshot2.png">

**Additional note:** The successful execution of this command means the restore process was initiated. Theprocess of restoring the backup takes some time to complete. For the convenience of the user, the restore process runs in the background. If you wish to know the progress and/or the state of the restore, you can use the &quot;cf service SERVICE\_INSTANCE\_NAME&quot; command.

### Aborting a restore:

**Command:** cf abort-restore SERVICE\_INSTANCE\_NAME

**Usage:** This command is used to abort the process of restore, of the service\_instance, previously initiatedby you. You need to provide the name of the service-instance as the parameter. Upon successful execution of the command, the plugin will initiate the aborting of the restore process.

**Expected Output:**

Aborting restore for [SERVICE\_INSTANCE\_NAME]

Success!

Restore has been aborted for the instance name: [SERVICE\_INSTANCE\_NAME]

**Example:**

<img src="https://github.com/SAP/service-fabrik-cli-plugin/blob/master/images/plugin_screenshot3.png">

**Additional note:** The successful execution of this command means the abort process was initiated. Theprocess of aborting the backup again takes some time to complete. For the convenience of the user, the abort process too runs in the background. If you wish to know the progress and/or the state of the backup, you can use the &quot;cf service SERVICE\_INSTANCE\_NAME&quot; command.

## [Error status](#error-status)

During the course of using the various commands of the plugin, you might come across various scenarios. Upon successful execution of any command, you get a message, &quot;Success!&quot; followed by the command specific information. But in case of an error, there are various status and error messages displayed to the user. Some of these erroneous status and their meaning are discussed here.

## **Starting a restore:**

### Unauthorized:

**Message** : Unauthorised/Access to resource forbidden

**Description** : This may have occurred because you don&#39;t have access to the space where the service-instance is or to the service instance itself. Please check if you have logged in the correct space or if you have typed in the correct instance-name or if you have the right permissions for the same.

### Another concurrent operation:

**Message** : Another operation is already in progress for the service instance

**Description** : You may be trying to apply a command on a service-instance which is already in process ofanother operation. The service-instance may be already undergoing a backup/restore operation. Kindly wait till this operation is over and then try again.

## **Aborting a restore:**

### No restore in progress:

**Message** : success (currently no restore in progress for this service instance)

**Description** : This isn&#39;t technically an error message. This occurs when you are trying to abort a restore ona service-instance whereas there are no restore-operations in process for this instance. The restore process may have already completed or no restore operation was made prior to.

### Unauthorized:

**Message** : Unauthorised

**Description** : This may have occurred because you don&#39;t have access to the space where the service-instance is or to the service instance itself. Please check if you have logged in the correct space or if you have typed in the correct instance-name or if you have the right permissions for the same.

# [Common Failure Scenarios](#common-failure-scenarios)

Apart from the erroneous status the user may encounter, as described in the previous section, the plugin handles some additional erroneous scenarios. These will be discussed in this section. We provide the error code, the commands where you might encounter such an error and the message you get when you encounter such an error. We also describe what the error actually means.

## Incorrect Number Of Arguments

**Triggered by:** User enters more arguments than permitted by the plugin/command.

The plugin supports a fixed set of arguments which you must provide while entering a command. e.g. If you enter more than one instance-name to take backup or don&#39;t provide an instance name while taking backup, this error can be triggered.

**Commands:** ALL

**Message:** You have entered an incorrect number of arguments. Enter &#39;cf backup&#39; to check the list ofcommands and their usage.

## Incorrect Space

**Triggered by:** User enters an instance name which doesn&#39;t belong to the targetted space.

The plugin supports the functionalities such that the user can use commands on service-instances present in the space targeted. This means the space in which you used &quot;cf-login&quot; to login to. Kindly ensure that in the operation you are trying to perform, you have mentioned a service-instance present in this space.

**Commands:** Start-Restore, Abort-Restore

**Message:  ** Instance  name  requested  doesn&#39;t  belong  to  the  org:  [ORG\_NAME]  and  the  space:

[SPACE\_NAME]. Please target the correct org and space.

## Incorrect Instance Name

**Triggered by** : User enters an instance name which doesn&#39;t exist.

Sometimes a user might have entered an incorrect instance name for an operation. That could trigger this error. Please check the spelling of the instance, if it exists or the case-sensitiveness of the instance.

**Commands:** Start-Restore, Abort-Restore

**Message:** Service Instance, &quot;Instance Name&quot;doesn&#39;t exist.

## Incorrect Command Usage

**Triggered by:** User enters any command not supported by the plugin.

The user might have entered a command not supported by the plugin or wrong spelling of a command.

Please enter &quot;cf help&quot; to see a list of commands supported by the plugin and check their spelling/usage.

**Commands:** NA

Message: This is not a registered command. See &#39;cf help&#39; (This is handled by CF itself.)

## User Logged Out Error

**Triggered by:** User attempts to run any command while being logged out of cf domain.

You may have been logged out of cf. Please enter &quot;cf login&quot; in your cli and login to continue.

**Commands:** ALL

**Message:** No Access Token was found. You may be logged out. Please log in to continue.

## Multiple Guid for a deleted service instance Error

**Triggered by:** User attempts to fetch backups for deleted service-insatance.

The deleted service instance maps to multiple GUIDs. 

**Commands:** cf list-backup [SERVICE\_INSTANCE\_NAME] --deleted

**Message:** [SERVICE\_INSTANCE\_NAME] maps to multiple instance GUIDs, please use 'cf instance-events --delete' to list all instance delete events, get required instance guid from the list and then use 'cf list-backup --guid GUID' to fetch backups list.
Enter 'cf backup' to check the list of commands and their usage.

## Deleted instance not found Error

**Triggered by:** User attempts to fetch backups for deleted service-instance.

The given deleted service instance does not exists. 

**Commands:** cf list-backup [SERVICE\_INSTANCE\_NAME] --deleted

**Message:** Instance Guid not found for the given deleted instance [SERVICE\_INSTANCE\_NAME].
Enter 'cf backup' to check the list of commands and their usage.
