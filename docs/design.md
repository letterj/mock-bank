# FCFC Mock Bank used for Testing
We will need a Mock Bank to Send and Recieve data on the 3 types of FCFC Bank transactions (custmer deposit, interest deposit and customer withdraw)

## Config Values
These values will be loaded into the api at spin up(there could be multiple currencies in the file)
**currencies**
* currency_code
* decimal_places
* active_saturday
* recon_time

## Log
Using the Apache log format: 

LogFormat "%h %l %u %t \"%r\" %>s %b" common.

breaking down what each section of that log means.
* %h - The IP address of the client.
* %l - The identity of the client determined by identd on the client’s machine. Will return a hyphen (-) if this information is not available.
* %u - The userid of the client if the request was authenticated.
* %t - The time that the request was received.
* \"%r\" - The request line that includes the HTTP method used, the requested resource path, and the HTTP protocol that the client used.
* %>s - The status code that the server sends back to the client.
* %b - The size of the object requested.

If a request was made to a website using the above-mentioned log format the resulting log would look similar to the following.


`127.0.0.1 - peter [9/Feb/2017:10:34:12 -0700] "GET /sample-image.png HTTP/2" 200 1479`

## Tables
Tables used in for persistent data storage.
**currency**
* currency_code
* decimal_places
* active_saturday
* recon_time

**customers**
 * ID
 * lei
 * Customer Name
 * QAccount

**accounts**
 * number
 * QAccount
 * currency_type
 * balance

**transactions**
* ID
* type  [deposit-customer, deposit-interest, withdrawl]
* date
* account_number
* customer_id
* description
* debit account
* credit account
* amount

**notifications**
* ID
* type [deposit withdraw]
* message_date
* customer_id
* account_id
* transaction_id
* message
* amount
* status
* ack


## Operations:
### Inital Setup
Action
* create tables (currencies, customers, accounts, transactions and notifications)
* create fcfc currency account(s)

### Customer
~~POST /customer~~
* lei
* Customer Name
* QAccount
* Currency
Action
* create customer record
    Returns Customer ID
* log operation

~~GET /customer~
Action:
* returns information on all customers
* log operation

~~GET /customer/{id}~~
Action:
* returns informtion about a single customer
* log operation

### Account
~~GET  /account~~
Action
* return information on all fcfc accounts at the bank
* log operation

~~GET /account/{AcctNumber}~~
* return information on a single account
* log operation

### Transactions
~~GET  /transaction~~
Action
* return information on all transactions
* log operation

~~GET /transaction/{account_ID}~~
Filters
* start_date
* end_date
* account_number
* customer_id 
* type
Action:
* return information on transactions for a specific account and filters
* log operation

### Deposits
POST /deposit
* type
* lei
* QAccount
* Currency
* Amount
Actions:
* create a deposit transaction
    if type = interest then lei and QAccount are blank
    if type = customer then lei and QAccount are required
* create a notification record
* returns transaction_id
* log operation

### Withdraws
POST /withdraws
* lei
* Currency
* Amount
* Home Bank Name
* Wire Instructions
Actions:
* create a withdraw transaction
* create a notification record
* returns transaction_id
* log operation

### Notifications
GET /notification
Filter:
* ack
* transaction_id"`
* account_number
* customer_id 
* start_date
* end_date
Actions:
* returns transaction_id
* log operation