# JCheck

## About this plugin
A Micro-UTP, plug-able sanity checker for any on-prem JFrog platform instance

## Installation with JFrog CLI
Installing the latest version:

`$ jfrog plugin install JCheck`

Installing a specific version:

`$ jfrog plugin install JCheck@version`

Uninstalling a plugin

`$ jfrog plugin uninstall JCheck`

## Usage
### Commands
* list
    - Arguments:
        - None
    - Flags:
        - None
    - Example:
    ```
  $ jfrog JCheck list

    Name                          Description                                                                   Group        Is Read Only
    RTConnectionCheck             Performs a check that validates that a connection to RT works                 Artifactory  true
    RTDeployCheck                 Deploy a large file to Artifactory, download and verify checksum              Artifactory  false
    RTHasRepositoriesCheck        Performs a check that validates that RT has configured repositories           Artifactory  true
    SelfCheck                     A sanity check that should pass                                               Self         true
    XrayConnectionCheck           Performs a check that validates that a connection to XRAY works               Xray         true
    XrayDbConnectionPoolCheck     Performs a check that DB connection pool is not maxed                         Xray         true
    XrayFreeDiskSpaceCheck        Performs a check that free disk space is above 100Gb                            Xray         true
    XrayHasIndexedResourcesCheck  Performs a check that validates that XRAY has configured indexed resources    Xray         true
    XrayHasPoliciesCheck          Performs a check that validates that XRAY has configured policies             Xray         true
    XrayHasWatchesCheck           Performs a check that validates that XRAY has configured watches              Xray         true
    XrayMonitoringAPICheck        Performs a check that calls XRAY monitoring API                               Xray         true
    XrayRabbitMQCheck             Performs a check that critical queues are not overflowed                      Xray         true
    XrayViolationCountCheck       Performs a check that checks that Xray is not generating too many violations  Xray         true
  ```

* check
    - Arguments:
        - what - The name (or partial name) of the checks you want to run. It can be a group name, check name or ALL for all
    - Flags:
        - readOnlyMode: Only run checks which are read only. **[Default: false]**
        - loop: Loop over times. **[Default: 1]**
        - loopSleep: Sleep time (in seconds) between loops. **[Default: 0]**
        - json: Return JSON result. **[Default: false]**
    - Example:
    ```
  $ jfrog JCheck ALL
  
    ** Running check: ...
    ......
    ......

    Name                          Failure Ind  Message
    RTConnectionCheck                          RT version 7.27.10 was detected
    RTDeployCheck                              Artifact deployed and downloaded successfully
    RTHasRepositoriesCheck                     detected 14 repositories
    SelfCheck                                  Self check passed
    XrayConnectionCheck                        Xray version 3.35.0 was detected
    XrayDbConnectionPoolCheck                  Xray DB connection pool has available connections (0/60 connections)
    XrayFreeDiskSpaceCheck        FAIL         Xray disk free space is lower than 100Gb (88.06 Gb)
    XrayHasIndexedResourcesCheck               detected 11 indexed repositories
    XrayHasPoliciesCheck                       detected 1 policies
    XrayHasWatchesCheck                        detected 1 watches
    XrayMonitoringAPICheck        FAIL         Server response: 403 Forbidden
    XrayRabbitMQCheck                          Total number of messages = 0
    XrayViolationCountCheck                    detected 11 violations in last 24 hours
    
    
  ```

### Environment variables
None

## Additional info
None.

## Release Notes
The release notes are available [here](RELEASE.md).
