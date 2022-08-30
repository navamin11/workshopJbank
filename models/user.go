package models

import (
	"context"
	"fmt"
	"time"
	"workshop/configs"
	"workshop/utils"

	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string    `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password,omitempty"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Useflag  string    `json:"useflag"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

type SessionLogin struct {
	AccessToken string    `json:"access_token"`
	AccessUuid  string    `json:"access_uuid"`
	UserId      string    `json:"user_id"`
	Name        string    `json:"name"`
	Expire      time.Time `json:"expire"`
}

var ctx = context.Background()

func CreateUser(username, password, name, email string) (User, error) {

	var user User

	if err := user.searchCreateUser(username); err != nil {
		return user, err
	}

	user.Id = uuid.NewV4().String()
	user.Username = utils.Trim(username)
	user.Password = utils.Hash(password)
	user.Name = name
	user.Email = email
	user.Useflag = "Y"
	user.CreateAt = utils.LocalTime(time.Now())
	user.UpdateAt = utils.LocalTime(time.Now())

	if _, err := user.insertUser(); err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) Profile() (*User, error) {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	if err := db.NewSelect().Model(u).Where("id = ?", u.Id).Scan(ctx); err != nil {
		return u, err
	}

	u.CreateAt = utils.LocalTime(u.CreateAt)
	u.UpdateAt = utils.LocalTime(u.UpdateAt)
	u.PrepareGive()

	return u, nil
}

func Login(username, password string) (SessionLogin, error) {
	var (
		session SessionLogin
		user    User
	)

	user, err := searchLoginUser(username, password)
	if err != nil {
		return session, err
	}

	token, err := utils.GenerateToken(user.Id, user.Name)
	if err != nil {
		return session, err
	}

	if err := createKey(token.AccessUuid, user.Id, time.Unix(token.AtExpires, 0)); err != nil {
		return session, err
	}

	session.AccessToken = token.AccessToken
	session.AccessUuid = token.AccessUuid
	session.UserId = user.Id
	session.Name = user.Name
	session.Expire = utils.LocalTime(time.Unix(token.AtExpires, 0))

	return session, nil
}

func Logout(token string) error {
	data, err := utils.ExtractTokenMetadataApi(token)
	if err != nil {
		return err
	}

	rdb := configs.ConnectRedis()
	defer rdb.Close()

	if err := rdb.Del(ctx, data.AccessUuid).Err(); err != nil {
		return err
	}

	return nil
}

func createKey(key, value string, t time.Time) error {
	rdb := configs.ConnectRedis()
	defer rdb.Close()

	if err := rdb.Set(ctx, key, value, t.Sub(time.Now())).Err(); err != nil {
		return err
	}
	return nil
}

func searchLoginUser(username, password string) (User, error) {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	var u User

	if err := db.NewSelect().Model(&u).Where("username = ? AND useflag = ?", username, "Y").Scan(ctx); err != nil {
		err := fmt.Errorf("Username not found")
		return u, err
	}

	if err := utils.VerifyPassword(password, u.Password); err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		err := fmt.Errorf("Password is incorrect")
		return u, err
	}
	return u, nil
}

func (u *User) searchCreateUser(username string) error {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	if err := db.NewSelect().Model(u).Where("username like ?", username).Column("id").Scan(ctx); err == nil {
		err := fmt.Errorf("Username not available")
		return err
	}
	return nil
}

func (u *User) insertUser() (*User, error) {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	if _, err := db.NewInsert().Model(u).Exec(ctx); err != nil {
		return u, err
	}

	u.PrepareGive()

	return u, nil
}

func (u *User) PrepareGive() {
	u.Password = ""
}
