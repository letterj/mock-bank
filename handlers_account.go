package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) getAccounts(w http.ResponseWriter, r *http.Request) {
	acts, err := getAccounts(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, acts)
}

func (a *App) getAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cur := vars["currency"]

	act := account{CurrencyCode: cur}
	if err := act.getAccount(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Account not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, act)
}
func (a *App) updateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acctNumber := vars["acctNumber"]

	var act account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&act); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Account Data")
		return
	}
	defer r.Body.Close()

	act.AcctNumber = acctNumber
	if err := act.updateAccount(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Account not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, act)

}
