// main_test.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

var a App

func TestMain(m *testing.M) {
	val, err := LoadConfig("default")

	a = App{}
	a.Initialize(val)

	rows, err := a.DB.Query("SELECT count() FROM currencies")
	if err != nil {
		log.Fatal(err)
	}

	var cnt int
	for rows.Next() {
		rows.Scan(&cnt)
	}

	code := m.Run()

	//clearTables()

	os.Exit(code)
}

func clearTables() {
	a.DB.Exec("DELETE FROM customers")
	a.DB.Exec("DELETE FROM transactions")
	a.DB.Exec("DELETE FROM notifications")
}

func TestVersion(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/version", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	r := message{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	fmt.Println("VERSION: ", r.Message)

	if r.Message == "" {
		t.Errorf("Expected %s currency(ies). Got %s", "Non-blank", r.Message)
	}
}

func TestCurrencyTable(t *testing.T) {
	numRows := 2

	req, _ := http.NewRequest("GET", "/api/v1/currency", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	r := []Currency{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d currency(ies). Got %d", numRows, len(r))
	}
}

func TestAccountTable(t *testing.T) {
	// clearTables()
	numRows := 2

	req, _ := http.NewRequest("GET", "/api/v1/account", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	r := []account{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d account(s). Got %d", numRows, len(r))
	}
}

func TestCustomerTable(t *testing.T) {
	clearTables()
	numRows := 0

	req, _ := http.NewRequest("GET", "/api/v1/customer", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	r := []customer{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d customer(s). Got %d", numRows, len(r))
	}
}

func TestTransactionTable(t *testing.T) {
	clearTables()
	numRows := 0

	req, _ := http.NewRequest("GET", "/api/v1/transaction", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	r := []transaction{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d transaction(s). Got %d", numRows, len(r))
	}
}

func TestNotificationTable(t *testing.T) {
	clearTables()
	numRows := 0

	req, _ := http.NewRequest("GET", "/api/v1/notification", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	r := []notification{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d currencies. Got %d", numRows, len(r))
	}
}

func TestSpecificAccount(t *testing.T) {
	clearTables()
	numRows := 0
	aNumber := "12345678-091"

	url := fmt.Sprintf("/api/v1/account/%s", aNumber)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", `application/json`)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	r := []account{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d account(s). Got %d", numRows, len(r))
	}
}

func TestSpecificCustomer(t *testing.T) {
	clearTables()
	numRows := 0
	anLei := "01-23456789"

	url := fmt.Sprintf("/api/v1/customer/%s", anLei)
	req, _ := http.NewRequest("GET", url, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	r := []customer{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d customer(s). Got %d", numRows, len(r))
	}
}

func TestSpecificTransaction(t *testing.T) {
	clearTables()
	numRows := 0
	transID := 1

	url := fmt.Sprintf("/api/v1/transaction/%d", transID)
	req, _ := http.NewRequest("GET", url, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	r := []transaction{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d transaction(s). Got %d", numRows, len(r))
	}
}

func TestSpecificNotification(t *testing.T) {
	clearTables()
	numRows := 0
	noticeID := 1

	url := fmt.Sprintf("/api/v1/transaction/%d", noticeID)
	req, _ := http.NewRequest("GET", url, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	r := []notification{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d notification(s). Got %d", numRows, len(r))
	}
}

func TestAddCustomer(t *testing.T) {
	clearTables()
	numRows := 1

	data := map[string]interface{}{
		"lei":            "123456-00",
		"name":           "Test Trading Co of America",
		"quorum_account": "0x111111",
	}

	bytesRepresentation, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	reqA, _ := http.NewRequest("POST", "/api/v1/customer", bytes.NewBuffer(bytesRepresentation))
	reqA.Header.Set("Content-Type", "application/json")
	presponse := executeRequest(reqA)

	checkResponseCode(t, http.StatusOK, presponse.Code)

	req, _ := http.NewRequest("GET", "/api/v1/customer", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	r := []customer{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d customers. Got %d", numRows, len(r))
	}
}

func TestSendDeposit(t *testing.T) {
	clearTables()

	// POST CUSTOMER
	cdata := map[string]interface{}{
		"lei":            "123456-00",
		"name":           "Test Trading Co of America",
		"quorum_account": "0x111111",
	}

	bytesRepresentation_cust, err := json.Marshal(cdata)
	if err != nil {
		log.Fatalln(err)
	}

	reqC, _ := http.NewRequest("POST", "/api/v1/customer", bytes.NewBuffer(bytesRepresentation_cust))
	reqC.Header.Set("Content-Type", "application/json")
	cresponse := executeRequest(reqC)

	checkResponseCode(t, http.StatusOK, cresponse.Code)

	// POST DEPOSIT
	data := map[string]interface{}{
		"type":           "WIRE",
		"name":           "Test Trading Co of America",
		"quorum_account": "Ox111111",
		"currency_code":  "USD",
		"amount":         1111.00,
	}

	bytesRepresentation, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	reqD, _ := http.NewRequest("POST", "/api/v1/deposit", bytes.NewBuffer(bytesRepresentation))
	reqD.Header.Set("Content-Type", "application/json")
	dresponse := executeRequest(reqD)

	checkResponseCode(t, http.StatusOK, dresponse.Code)

	rd := []deposit{}
	rdbody, _ := ioutil.ReadAll(dresponse.Body)
	json.Unmarshal(rdbody, &rd)

	req, _ := http.NewRequest("GET", "/api/v1/transaction", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	r := []transaction{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if r[0].Amount != 1111.00 {
		t.Errorf("Expected transaction amount %v. Got %v", 1111.00, r[0].AcctNumber)
	}
}

func TestAllAccounts(t *testing.T) {
	clearTables()
	numRows := 2

	url := fmt.Sprintf("/api/v1/account")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", `application/json`)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	r := []account{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if len(r) != numRows {
		t.Errorf("Expected %d accounts. Got %d", numRows, len(r))
	}
}

func TestSpecificAccounts(t *testing.T) {
	clearTables()
	lookupCurrency := "USD"

	url := fmt.Sprintf("/api/v1/account/%s", lookupCurrency)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", `application/json`)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var r account
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)

	if r.CurrencyCode != lookupCurrency {
		t.Errorf("Expected account %s. Got %s", lookupCurrency, r.CurrencyCode)
	}
}

func TestSendWithdraw(t *testing.T) {
	clearTables()

	// POST CUSTOMER
	cdata := map[string]interface{}{
		"lei":            "123456-00",
		"name":           "Customer Withdraw",
		"quorum_account": "0x111111",
	}

	bytesRepresentation_cust, err := json.Marshal(cdata)
	if err != nil {
		log.Fatalln(err)
	}

	reqC, _ := http.NewRequest("POST", "/api/v1/customer", bytes.NewBuffer(bytesRepresentation_cust))
	reqC.Header.Set("Content-Type", "application/json")
	cresponse := executeRequest(reqC)

	checkResponseCode(t, http.StatusOK, cresponse.Code)

	// POST DEPOSIT
	pdata := map[string]interface{}{
		"type":           "WIRE",
		"name":           "Customer Withdraw",
		"quorum_account": "Ox111111",
		"currency_code":  "USD",
		"amount":         100000.00,
	}

	pbytesRepresentation, err := json.Marshal(pdata)
	if err != nil {
		log.Fatalln(err)
	}

	reqD, _ := http.NewRequest("POST", "/api/v1/deposit", bytes.NewBuffer(pbytesRepresentation))
	reqD.Header.Set("Content-Type", "application/json")
	dresponse := executeRequest(reqD)

	checkResponseCode(t, http.StatusOK, dresponse.Code)

	// ACCOUNT NUMBER
	aurl := fmt.Sprintf("/api/v1/account/%s", "USD")
	areq, _ := http.NewRequest("GET", aurl, nil)
	areq.Header.Add("Content-Type", `application/json`)
	aresponse := executeRequest(areq)

	checkResponseCode(t, http.StatusOK, aresponse.Code)

	var ra account
	abody, _ := ioutil.ReadAll(aresponse.Body)
	json.Unmarshal(abody, &ra)

	// POST WITHDRAW
	wdata := map[string]interface{}{
		"lei":            "123456-00",
		"account_number": ra.AcctNumber,
		"bank_name":      "First State Bank",
		"currency_code":  "USD",
		"amount":         100.00,
		"instructions":   "FOB 1234",
		"notes":          "Test Note",
	}

	bytesRepresentation, err := json.Marshal(wdata)
	if err != nil {
		log.Fatalln(err)
	}

	reqW, _ := http.NewRequest("POST", "/api/v1/withdraw", bytes.NewBuffer(bytesRepresentation))
	reqD.Header.Set("Content-Type", "application/json")
	wresponse := executeRequest(reqW)

	checkResponseCode(t, http.StatusOK, wresponse.Code)

	var wd withdraw
	body, _ := ioutil.ReadAll(wresponse.Body)
	json.Unmarshal(body, &wd)

	if wd.RefID == 0 {
		t.Errorf("Expected a reference id %v. Got %v", "Non-Blank", wd.RefID)
	}
}
