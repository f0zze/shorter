package handlers

import (
	json "encoding/json"
	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
	"net/http"
)

type UserHandler struct {
	Storage storage.Storage
	Service services.ShortURLService
}

func (u *UserHandler) Urls(resp http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value("userID").(string)

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
