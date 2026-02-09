package main

import (
	"encoding/json"
	"net/http"

	"github.com/boatnoah/iupload/internal/auth"
	"github.com/boatnoah/iupload/internal/storage"
)

func (a *app) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var userPayload auth.UserPayload

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&userPayload)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	_, session, err := a.auth.RegisterUser(r.Context(), userPayload)
	if err != nil {
		if err == storage.ErrorDuplicateUserName {
			http.Error(w, "Duplicate Error", http.StatusConflict)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sessionJson, err := json.Marshal(session)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Write(sessionJson)

}

func (a *app) logInUserHandler(w http.ResponseWriter, r *http.Request) {

	var userLoginPayload auth.UserLoginPayload

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&userLoginPayload)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := a.auth.LogInUser(r.Context(), userLoginPayload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sessionJson, err := json.Marshal(session)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Write(sessionJson)

}
