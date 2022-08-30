package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"workshop/models"
	"workshop/utils"

	"github.com/go-playground/validator/v10"
)

type RegisterInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// var validate *validator.Validate

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	var (
		input    RegisterInput
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

	data, err := models.CreateUser(input.Username, input.Password, input.Name, input.Email)
	if err != nil {
		utils.RespJSON(w, http.StatusBadRequest, models.Response{Code: "400", Message: err.Error()})
		return
	}

	utils.RespJSON(w, http.StatusOK, models.Response{Code: "200", Message: "success", Results: data})
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	var (
		input    LoginInput
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

	token, err := models.Login(input.Username, input.Password)
	if err != nil {
		utils.RespJSON(w, http.StatusUnauthorized, models.Response{Code: "401", Message: err.Error()})
		return
	}

	utils.RespJSON(w, http.StatusOK, models.Response{Code: "200", Message: "success", Results: token})
}

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	data, err := utils.ExtractTokenMetadataApi(token)
	if err != nil {
		utils.RespJSON(w, http.StatusUnauthorized, models.Response{Code: "401", Message: err.Error()})
		return
	}

	var user models.User

	user.Id = data.UserId

	resp, err := user.Profile()
	if err != nil {
		utils.RespJSON(w, http.StatusBadRequest, models.Response{Code: "400", Message: err.Error()})
		return
	}

	utils.RespJSON(w, http.StatusOK, models.Response{Code: "200", Message: "success", Results: resp})
}

func LogoutDelHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if err := models.Logout(token); err != nil {
		utils.RespJSON(w, http.StatusBadRequest, models.Response{Code: "400", Message: err.Error()})
		return
	}
	utils.RespJSON(w, http.StatusOK, models.Response{Code: "200", Message: "success"})
}
