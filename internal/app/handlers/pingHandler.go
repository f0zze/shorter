package handlers

import (
	"github.com/f0zze/shorter/internal/app/storage"
	"net/http"
)

type PingHandler struct {
	Storage storage.Storage
}

func (h *PingHandler) Get(resp http.ResponseWriter, req *http.Request) {
	hasConnection := h.Storage.Ping()

	if hasConnection {
		resp.WriteHeader(http.StatusOK)
		return
	}

	resp.WriteHeader(http.StatusInternalServerError)
}
