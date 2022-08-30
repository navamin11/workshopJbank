package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"workshop/models"
	"workshop/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type TransferInInput struct {
	AccountNo   string  `json:"account_no" validate:"required"`
	AccountName string  `json:"account_name" validate:"required"`
	Bank        string  `json:"bank" validate:"required"`
	DepositId   string  `json:"deposit_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
}

type TransferOutInput struct {
	AccountNo   string  `json:"account_no" validate:"required"`
	AccountName string  `json:"account_name" validate:"required"`
	Bank        string  `json:"bank" validate:"required"`
	DepositId   string  `json:"deposit_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
}

func TransferOutPostHandler(w http.ResponseWriter, r *http.Request) {
	var (
		input    TransferOutInput
		validate *validator.Validate
	)

	validate = validator.New()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespJSON(w, http.StatusUnprocessableEntity, models.Response{Code: "422", Message: err.Error()})
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		utils.RespJSON(w, http.StatusUnprocessableEntity, models.Response{Code: "422", Message: err.Error()})
		return
	}

	if err := validate.Struct(input); err != nil {
		utils.RespJSON(w, http.StatusBadRequest, models.Response{Code: "400", Message: err.Error()})
		return
	}

	var trans models.Transaction

	trans.AccountNo = input.AccountNo
	trans.AccountName = input.AccountName
	trans.Bank = input.Bank
	trans.Type = chi.URLParam(r, "type")
	trans.DepositId = input.DepositId
	trans.MoneyOut = input.Amount

	if err := trans.TransactionOut(); err != nil {
		utils.RespJSON(w, http.StatusBadRequest, models.Response{Code: "400", Message: err.Error()})
		return
	}

	utils.RespJSON(w, http.StatusOK, models.Response{Code: "200", Message: "success", Results: trans})
}

func TransferInPostHandler(w http.ResponseWriter, r *http.Request) {
	var (
		input    TransferInInput
		validate *validator.Validate
	)

	validate = validator.New()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespJSON(w, http.StatusUnprocessableEntity, models.Response{Code: "422", Message: err.Error()})
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		utils.RespJSON(w, http.StatusUnprocessableEntity, models.Response{Code: "422", Message: err.Error()})
		return
	}

	if err := validate.Struct(input); err != nil {
		utils.RespJSON(w, http.StatusBadRequest, models.Response{Code: "400", Message: err.Error()})
		return
	}

	var trans models.Transaction

	trans.AccountNo = input.AccountNo
	trans.AccountName = input.AccountName
	trans.Bank = input.Bank
	trans.Type = chi.URLParam(r, "type")
	trans.DepositId = input.DepositId
	trans.MoneyIn = input.Amount

	if err := trans.TransactionIn(); err != nil {
		utils.RespJSON(w, http.StatusBadRequest, models.Response{Code: "400", Message: err.Error()})
		return
	}

	utils.RespJSON(w, http.StatusOK, models.Response{Code: "200", Message: "success", Results: trans})
}
