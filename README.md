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
cd <repo_directory>/fc2-mock-bank 
make build 
```

## General Usage
```:bash
$ ./fc2-mock-bank.macos -help
Usage of ./fc2-mock-bank.macos:
  -c string
    	(short-hand) configuration file containing setup information (default "default")
  -config string
    	configuration file containing setup information (default "default")
  -v	(short-hand) application version
  -version
    	application version
```


## QuickStart
```:bash
make all 
``` 

## API Commands
### Health Check and Code Version
Code Version
```:bash
GET     /api/v1/version

{
    "message": "v0.2-beta"
}
```

### Accounts
List all Accounts
```:bash
GET     /api/v1/account

[
    {
        "acct_number": "98498081-00",
        "quorum_account": "0x",
        "currency_code": "USD",
        "balance": 0
    },
    {
        "acct_number": "19727887-01",
        "quorum_account": "0x",
        "currency_code": "EUR",
        "balance": 0
    }
]
```

Specific Account by Currency
```:bash
GET     /api/v1/account/{currency}  [USD|EUR|CAD|AUS|GBP|JPY]

[
      {
        "acct_number": "19727887-01",
        "quorum_account": "0x",
        "currency_code": "EUR",
        "balance": 0
    }
]
```

Update the Quorum Account number
```:bash
PUT     /api/vi/notification/{accout number}  [19727887-01]
-d {"quorum_account": "0x1234567898765432"}

{
    "acct_number": "19727887-01",
    "quorum_account": "0x1234567898765432",
    "currency_code": "EUR",
    "balance": 0
}
```


### Customers
List all Customers
```:bash
GET     /api/v1/customer

[
    {
        "id": 1,
        "lei": "123456-00",
        "name": "Test Customer 1",
        "quorum_account": "0x111111"
    }
]
```

Add a Customer
```:bash
POST    /api/v1/customer
-h "Content-Type: application/json"
-d {"lei": "123456-00", "name": "Test Customer 1", "quorum_account": "0x111111"}

{
    "id": 1,
    "lei": "123456-00",
    "name": "Test Customer 1",
    "quorum_account": "0x111111"
}
```


### Deposits
Add a Deposit
```:bash
POST    /api/v1/deposit
-h "Content-Type: application/json"
-d {"type": "WIRE","name": "Test Customer 1", "quorum_account": "Ox111111", "currency_code": "USD", "amount": 5000.50, "start_date": "", "end_date": "", "rate": 0}

{
    "type": "WIRE",
    "name": "Test Customer 1",
    "quorum_account": "Ox111111",
    "currency_code": "USD",
    "amount": 5000.5,
    "start_date": "",
    "end_date": "",
    "rate": 0,
    "refid": 1
}


POST    /api/v1/deposit
-h "Content-Type: application/json"
-d {"type": "INTEREST","name": "", "quorum_account": "", "currency_code": "USD", "amount": 19.63, "start_date":"01/01/2019", "end_date": "01/31/2019", "rate": 0.0112}

{
    "type": "INTEREST",
    "name": "",
    "quorum_account": "",
    "currency_code": "USD",
    "amount": 19.63,
    "start_date": "01/01/2019",
    "end_date": "01/31/2019",
    "rate": 0.0112,
    "refid": 2
}
```	


### Withdrawl
Submit a Withdrawl Request
```:bash
POST    /api/v1/withdraw
-h "Content-Type: application/json"
-d {"lei": "123456-00","account_number": "98498081-00", "bank_name": "First State Bank", "currency_code": "USD", "amount": 100.00,"instructions": "FOB 1234", "notes": "Test Note"}

{
    "lei": "123456-00",
    "account_number": "98498081-00",
    "currency_code": "USD",
    "amount": 100,
    "bank_name": "First State Bank",
    "instructions": "FOB 1234",
    "refid": 3,
    "notes": "Test Note"
}

```	


### Transactions
List all Transactions
```:bash
GET     /api/v1/transaction

[
    {
        "id": 1,
        "trans_type": "WIRE",
        "currency": "USD",
        "trans_date": "2019-07-17T09:28:32Z",
        "account_number": "98498081-00",
        "customer_id": 1,
        "quorum_account": "Ox111111",
        "description": "Customer Deposit",
        "amount": 5000.5,
        "start_date": "",
        "end_date": "",
        "rate": 0,
        "status": "POSTED"
    },
    {
        "id": 2,
        "trans_type": "INTEREST",
        "currency": "USD",
        "trans_date": "2019-07-17T09:30:00Z",
        "account_number": "98498081-00",
        "customer_id": 0,
        "quorum_account": "",
        "description": "Interest Rate Deposit",
        "amount": 19.63,
        "start_date": "01/01/2019",
        "end_date": "01/31/2019",
        "rate": 0.0112,
        "status": "POSTED"
    },
    {
        "id": 3,
        "trans_type": "WITHDRAW",
        "currency": "USD",
        "trans_date": "2019-07-17T09:34:10Z",
        "account_number": "98498081-00",
        "customer_id": 1,
        "quorum_account": "",
        "description": "FOB 1234",
        "amount": -100,
        "start_date": "",
        "end_date": "",
        "rate": 0,
        "status": "POSTED"
    }
]
```

### Notifications
List all Notifications
```:bash
GET     /api/v1/notification

[
    {
        "id": 1,
        "type": "WIRE",
        "notice_date": "2019-07-17T09:28:32Z",
        "currency": "USD",
        "customer_id": 1,
        "transaction_id": 1,
        "message": "Customer Deposit",
        "amount": 5000.5,
        "quorum_account": "Ox111111",
        "start_date": "",
        "end_date": "",
        "rate": 0,
        "status": "POSTED",
        "ack": false
    },
    {
        "id": 2,
        "type": "INTEREST",
        "notice_date": "2019-07-17T09:30:00Z",
        "currency": "USD",
        "customer_id": 0,
        "transaction_id": 2,
        "message": "Interest Rate Deposit",
        "amount": 19.63,
        "quorum_account": "",
        "start_date": "01/01/2019",
        "end_date": "01/31/2019",
        "rate": 0.0112,
        "status": "POSTED",
        "ack": false
    },
    {
        "id": 3,
        "type": "WITHDRAW",
        "notice_date": "2019-07-17T09:34:10Z",
        "currency": "USD",
        "customer_id": 1,
        "transaction_id": 3,
        "message": "FOB 1234",
        "amount": -100,
        "quorum_account": "",
        "start_date": "",
        "end_date": "",
        "rate": 0,
        "status": "POSTED",
        "ack": false
    }
]
```
Acknowledge a Notification and remove it from future GET Requests
```:bash
PUT     /api/vi/notification/{ID}  [2]

{
    "id": 2,
    "type": "INTEREST",
    "notice_date": "2019-07-17T09:30:00Z",
    "currency": "USD",
    "customer_id": 0,
    "transaction_id": 2,
    "message": "Interest Rate Deposit",
    "amount": 19.63,
    "quorum_account": "",
    "start_date": "01/01/2019",
    "end_date": "01/31/2019",
    "rate": 0.0112,
    "status": "POSTED",
    "ack": true
}
```

