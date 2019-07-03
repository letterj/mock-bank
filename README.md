# Setup for fcfcmockbank
A bank API used for testing the interacts between the 9th Gear platform an a Custodial bank


##  Install golang
Install (golang)[https://golang.org/]


##  Add in dependencies
Adding in the dependencies

```:bash
go get -v github.com/mattn/go-sqlite3
go get -v github.com/gorilla/mux
go get -v github.com/gorilla/handler
```

## Create an sqlite3 bank Database
The API works with a sqlite3 database called **bank.db**. This database is created when the API is launched and con.  


## Configuration file
The file **config.json.sample** has the appropiate values.  


##  Build the api
This will build the go application into an executable

```:bash
cd <repo_directory>/fcfcmockbank 
make build 
```

## Run the api in dev mode
```:bash
./fcfcmockbank [default | config.json]
```

## QuickStart
```:bash
make all 
``` 

## API Commands

### Accounts
List all accounts
```:bash
/api/v1/account
```

### Customers
List all Customers
```:bash
GET     /api/v1/customer
```
Add a Customer
```:bash
POST    /api/v1/customer
-h "Content-Type: application/json"
-d {"lei": "123456-00", "name": "Test Customer 1", "quorum_account": "0x111111"}
```

### Transactions
List all Transactions
```:bash
GET     /api/v1/transaction
```

### Notifications
List all Notifications
```:bash
GET     /api/v1/notification
```

### Deposits
Add a Deposit
```:bash
POST    /api/v1/deposit
-h "Content-Type: application/json"
-d {"type": "WIRE","name": "Test Customer 1", "quorum_account": "Ox111111", "currency_code": "USD", "amount": 1111.00}
```	
