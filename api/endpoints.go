package api

import (
	"encoding/json"
	"fp-dim-aws-guard-duty-ingress/internal"
	"fp-dim-aws-guard-duty-ingress/internal/config"
	"fp-dim-aws-guard-duty-ingress/internal/structs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type HttpResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HandleIncomingData(w http.ResponseWriter, r *http.Request) {
	urlToken := r.URL.Query().Get("token")

	if urlToken == "" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(&HttpResponse{
			Status:  http.StatusForbidden,
			Message: "Tokens not supplied",
		})
	}

	savedToken := viper.GetString("url-token")

	if urlToken != savedToken {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&HttpResponse{
			Status:  http.StatusUnauthorized,
			Message: "Tokens do not match",
		})
		return
	}

	item := structs.IncomingItem{}

	err := json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(&HttpResponse{
			Status:  http.StatusNotAcceptable,
			Message: "could not decode json into entity",
		})
		return
	}

	resp, status, err := internal.PushDataToController(item)

	if err != nil {
		logrus.Error(err)
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(*resp)
}

func SendHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ConfigEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(config.ReadConfig())
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
