package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/boatnoah/iupload/internal/storage"
)

type userKey string

const userCtx userKey = "user"

func (a *app) uploadImageHandler(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("image")

	if err != nil {
		http.Error(w, "image file is required", http.StatusBadRequest)
		return
	}

	defer file.Close()

	contentType := header.Header.Get("Content-Type")

	if contentType != "image/jpeg" && contentType != "image/png" {
		http.Error(w, "only JPG and PNG are allowed", http.StatusBadRequest)
		return
	}

	user, _ := r.Context().Value(userCtx).(*storage.User)

	image, err := a.svc.UploadImage(r.Context(), user.ID, header.Filename, file, contentType)

	if err != nil {
		http.Error(w, "Unable to upload image", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	imageJson, err := json.Marshal(image)

	if err != nil {
		http.Error(w, "unable to marshal image", http.StatusInternalServerError)
		return
	}

	w.Write(imageJson)

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
