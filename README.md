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

