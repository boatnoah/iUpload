package main

import (
	"encoding/json"
	"net/http"

	"github.com/boatnoah/iupload/internal/storage"
)

type userKey string

const userCtx userKey = "user"

func (a *app) uploadImageHandler(w http.ResponseWriter, r *http.Request) {

	user, _ := r.Context().Value(userCtx).(*storage.User)

	err := a.svc.UploadImage(r.Context(), user.ID, "bob", r.Body, "jpeg")
	if err != nil {
		http.Error(w, "Unable to upload image", http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "unable marshal", http.StatusInternalServerError)
		return
	}

	w.Write(userJson)

}

func (a *app) getImageByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))

}

func (a *app) transformImageHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("hello"))

}

func (a *app) deleteImageHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("hello"))
}
