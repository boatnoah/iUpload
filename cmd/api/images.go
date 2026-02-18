package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/boatnoah/iupload/internal/processor"
	"github.com/boatnoah/iupload/internal/storage"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
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

	var imagePayload processor.ImagePayload

	imagePayload.UserID = user.ID
	imagePayload.FileName = header.Filename
	imagePayload.Body = file
	imagePayload.ContentType = contentType

	image, err := a.svc.UploadImage(r.Context(), imagePayload)

	if err != nil {
		http.Error(w, "Unable to upload image", http.StatusInternalServerError)
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
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "id is not of type UUID from", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	reader, contentType, err := a.svc.GetByImageId(ctx, id)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	defer reader.Close()

	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Unable to send image", http.StatusInternalServerError)
		return
	}

}

func (a *app) deleteImageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "id is not of type UUID from", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err = a.svc.DeleteByImageId(ctx, id)
	if err != nil {

		if err == processor.ErrorNotFound {
			http.Error(w, "Could not find image", http.StatusNotFound)
			return
		}

		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *app) transformImageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "id is not in UUID form", http.StatusBadRequest)
		return
	}

	var tranformPayload processor.OperationPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&tranformPayload)

	if err != nil {
		http.Error(w, "Unable decode json", http.StatusBadRequest)
		return
	}

	if tranformPayload.Transformation == nil {
		http.Error(w, "Must have some fields", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	bytes, err := a.svc.TranformImage(ctx, id, tranformPayload)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to transform image", http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}
