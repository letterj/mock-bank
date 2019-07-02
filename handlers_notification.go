package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) getNotifications(w http.ResponseWriter, r *http.Request) {
	nots, err := getNotifications(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nots)
}

func (a *App) getNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Notification ID")
		return
	}

	n := notification{ID: id}
	if err := n.getNotification(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Notification not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, n)
}
