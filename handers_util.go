package main

import (
	"encoding/json"
	"net/http"
)

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
	}

	depData.RefID, err = depData.postDeposit(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, depData)
}
