package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
)

type UserHandler struct {
	Storage storage.Storage
	Service services.ShortURLService
}

func (u *UserHandler) Urls(resp http.ResponseWriter, req *http.Request) {
	userToken, err := req.Cookie("ID")

	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := services.GetUserID(userToken.Value)

	if userID == "" {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	list, _ := u.Service.FindByUser(userID)

	if list == nil {
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	jsonList, err := json.Marshal(list)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(jsonList)
}

func (u *UserHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	userToken, err := r.Cookie("ID")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := services.GetUserID(userToken.Value)

	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var urls []string

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&urls); err != nil {
		http.Error(w, "Could not parse body ", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	err = u.Service.DeleteURL(urls, userID)

	if err != nil {
		http.Error(w, "Could not delete urls", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
