package main

import (
	"database/sql"
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
	id := vars["id"]

	act := account{AcctNumber: id}
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
