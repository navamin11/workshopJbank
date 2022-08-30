package controllers

import (
	"net/http"
	"workshop/models"
	"workshop/utils"
)

func DepositGetHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	data, err := utils.ExtractTokenMetadataApi(token)
	if err != nil {
		utils.RespJSON(w, http.StatusUnauthorized, models.Response{Code: "401", Message: err.Error()})
		return
	}

	var account models.Account

	account.UserId = data.UserId
	account.Flag = "main"

	resp, err := account.ListAcount()
	if err != nil {
		utils.RespJSON(w, http.StatusBadRequest, models.Response{Code: "400", Message: err.Error()})
		return
	}

	utils.RespJSON(w, http.StatusOK, models.Response{Code: "200", Message: "success", Results: resp})
}

func DepositListGetHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	data, err := utils.ExtractTokenMetadataApi(token)
	if err != nil {
		utils.RespJSON(w, http.StatusUnauthorized, models.Response{Code: "401", Message: err.Error()})
		return
	}

	var account models.Account

	account.UserId = data.UserId

	resp, err := account.ListAcounts()
	if err != nil {
		utils.RespJSON(w, http.StatusBadRequest, models.Response{Code: "400", Message: err.Error()})
		return
	}

	utils.RespJSON(w, http.StatusOK, models.Response{Code: "200", Message: "success", Results: resp})
}
