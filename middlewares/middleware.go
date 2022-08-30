package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"workshop/configs"
	"workshop/models"
	"workshop/utils"
)

var ctx = context.Background()

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if token == "" {
			utils.RespJSON(w, http.StatusUnauthorized, models.Response{Code: "401", Message: "Missing Authorization Header"})
			return
		} else {
			if err := utils.TokenValidApi(token); err != nil {
				utils.RespJSON(w, http.StatusUnauthorized, models.Response{Code: "401", Message: err.Error()})
				return
			}

			if err := checkUseflag(token); err != nil {
				utils.RespJSON(w, http.StatusUnauthorized, models.Response{Code: "401", Message: err.Error()})
				return
			}

			if err := checkTokenExpire(token); err != nil {
				utils.RespJSON(w, http.StatusUnauthorized, models.Response{Code: "401", Message: err.Error()})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func checkUseflag(token string) error {
	data, err := utils.ExtractTokenMetadataApi(token)
	if err != nil {
		return err
	}

	db := configs.ConnectPostgreSQL()
	defer db.Close()

	var u models.User

	if err := db.NewSelect().Model(&u).Where("id = ? AND useflag = ?", data.UserId, "Y").Column("id", "name", "useflag").Scan(ctx); err != nil {
		err := fmt.Errorf("User is not authorized to access this application")
		return err
	}

	return nil
}

func checkTokenExpire(token string) error {
	data, err := utils.ExtractTokenMetadataApi(token)
	if err != nil {
		return err
	}

	rdb := configs.ConnectRedis()
	defer rdb.Close()

	tokenKey, err := rdb.Get(ctx, data.AccessUuid).Result()
	if err != nil || tokenKey == "" {
		err := fmt.Errorf("Authorization Token not found")
		return err
	}

	return nil
}
