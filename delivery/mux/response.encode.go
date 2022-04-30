package mux

import (
	"encoding/json"
	"net/http"

	"github.com/aditya37/geofence-service/util"
)

type ErrorEncoder interface {
	Err() error
}

func EncodeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
	return
}

func EncodeErrorResponse(w http.ResponseWriter, data interface{}) {
	err := data.(*util.ErrorMsg)
	w.WriteHeader(err.HttpRespCode)
	json.NewEncoder(w).Encode(err)
	return
}
