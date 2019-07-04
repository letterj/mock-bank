package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) getTransactions(w http.ResponseWriter, r *http.Request) {
	trans, err := getTransactions(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, trans)
}

func (a *App) getTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Transaction ID")
		return
	}

	tran := transaction{ID: id}
	if err := tran.getTransaction(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Transaction not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, tran)
}
