
# Release Notes and WIP

## v0.4-beta
Full Deposit and Withdraw with TLS support

New:
* Added TLS support 

Modifications:
* broke out code into separate files for cleaner layout


## v0.3.1-beta
Full Deposit and Withdraw 

Bug:
* Notification ack was not being set properly
  * for Withdrawls  ack = true
  * for Deposits    ack = false

## v0.3-beta
Full Deposit and Withdraw

New:
* general request logging
* Added GET /version to return the current code version
* command line "help"
* using golang package "flags" to work with command line
* release_notes.md
* Added fields to the notification and transaction tables
  * Quorum Account number
  * Start Date
  * End Date
  * Rate
* Added PUT on accounts to be able to update Quorum Account
* Added unique constraints on table columns
  * accounts.qaccount
  * accunts.currency_code
  * customer.qaccount
  * customer.lei

Modifications:
* converted fmt.Print statements to logs
* command line usage presentation
* updated README.md file to include result examples
* command line default
* broke up deposit and withdraw validation into sub functions 
* default quorum accounts created for accounts were realistic values
* default fc2 account number were realistic
* objects returns or POSTs and PUTs had appropriate values
  * Notification
  * Account
  * Customer
* cleaned up golang tests and verified they all pass


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