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
