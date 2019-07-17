package main

import (
	"encoding/json"
	"net/http"
)

func (a *App) version(w http.ResponseWriter, r *http.Request) {
	msg := message{Message: version}
	respondWithJSON(w, http.StatusOK, msg)
}

func (a *App) getCurrencies(w http.ResponseWriter, r *http.Request) {
	curs, err := getCurrencies(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, curs)
}

func (a *App) addDeposit(w http.ResponseWriter, r *http.Request) {
	var err error

	var depData deposit
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&depData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Deposit Values")
		return
	}
	defer r.Body.Close()

	msg := validateDeposit(a.DB, depData)
	if msg != "" {
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	err = depData.postDeposit(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, depData)
}

func (a *App) addWithdraw(w http.ResponseWriter, r *http.Request) {
	var err error

	var wData withdraw
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&wData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Withdraw Values")
		return
	}
	defer r.Body.Close()

	msg := validateWithdraw(a.DB, wData)
	if msg != "" {
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	err = wData.postWithdraw(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, wData)
}
