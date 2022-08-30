package utils

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// respondwithJSON write json response format
func RespJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Hash(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func Trim(username string) string {
	html.EscapeString(strings.TrimSpace(strings.ReplaceAll(username, " ", "")))
	return username
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LocalTime(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return t.In(loc)
}
