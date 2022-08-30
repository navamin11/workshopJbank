package controllers

import (
	"net/http"
	"workshop/models"
	"workshop/utils"
)

func HealthCheckGetHandler(w http.ResponseWriter, r *http.Request) {

	res := make(map[string]interface{})

	res["msg"] = "Welcome to JBank"

	utils.RespJSON(w, http.StatusOK, models.Response{Code: "200", Message: "success", Results: res})
}
