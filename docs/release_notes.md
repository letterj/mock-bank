
# Release Notes and WIP

## v0.3-beta
Title

New:
* general request logging
* ability to determine what version of code is running
* command line "help"
* using golang package "flags" to work with command line
* release_notes.md

Modifications:
* converted fmt.Print statements to logs
* command line usage presentation
* updated README.md file
* command line default 

Bugs:


 
## v0.2-beta
~~Upgrading Notification Functionality and Bug Fixes~~

Bug Fixes:
* Invalid deposit data does not create a transaction
* Fixed spelling of the word currency
* Fixed spelling of the word deposit

Functionality Changes:
* GET notifications will now only return notifications that have the "ack" field set to FALSE

Additional Functionality:
* PUT on a specific notification will set "ack" to TRUE and it will not longer be returned on a GET


## v0.1-beta
~~This is the initial release of the fc2-mock-bank. It's intended for use in testing only~~
* Working through the process 