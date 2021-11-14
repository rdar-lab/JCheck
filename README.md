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

  Name        Description   Group       Is Read Only
  Check...    Desc...       Group....   true/false  
  Check...    Desc...       Group....   true/false  
  Check...    Desc...       Group....   true/false  

  ```

* check
    - Arguments:
        - what - The names of the checks you want to run. It can be a group name, check name or ALL for all
    - Flags:
        - readOnlyMode: Only run checks which are read only. **[Default: false]**
        - loop: Loop over times. **[Default: 1]**
    - Example:
    ```
  $ jfrog JCheck ALL
  
    ** Running check: check1...
    Finished running check: check1, result=true, message=Everything OK

    Name        Is Success  Message
    Check...    true        Everything OK
    Check...    true        Everything OK
    Check...    true        Everything OK
     
  ```

### Environment variables
None

## Additional info
None.

## Release Notes
The release notes are available [here](RELEASE.md).
