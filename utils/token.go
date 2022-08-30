package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/twinj/uuid"
)

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
	RtExpires   int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     string
	Name       string
}

func GenerateToken(userid, name string) (*TokenDetails, error) {

	var err error

	lifespan, err := strconv.Atoi(os.Getenv("MINUTE_LIFESPAN"))
	if err != nil {
		return nil, err
	}

	var td TokenDetails

	td.AtExpires = time.Now().Add(time.Minute * time.Duration(lifespan)).Unix()
	td.AccessUuid = uuid.NewV4().String()

	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["userid"] = userid
	atClaims["name"] = name
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_KEY")))
	if err != nil {
		return nil, err
	}

	return &td, nil
}

func ExtractToken(bearerToken string) string {
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func VerifyTokenApi(bearerToken string) (*jwt.Token, error) {
	tokenString := ExtractToken(bearerToken)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func TokenValidApi(bearerToken string) error {
	token, err := VerifyTokenApi(bearerToken)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadataApi(bearerToken string) (AccessDetails, error) {

	var ad AccessDetails

	token, err := VerifyTokenApi(bearerToken)
	if err != nil {
		return ad, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return ad, err
		}

		userId, ok := claims["userid"].(string)
		if !ok {
			return ad, err
		}

		name, ok := claims["name"].(string)
		if !ok {
			return ad, err
		}

		ad.AccessUuid = accessUuid
		ad.UserId = userId
		ad.Name = name

		return ad, nil
	}

	return ad, err
}
