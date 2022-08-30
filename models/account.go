package models

import (
	"time"
	"workshop/configs"
	"workshop/utils"

	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"table:deposits"`
	Id            string    `json:"id"`
	AccountName   string    `json:"account_name"`
	AccountNo     string    `json:"account_no"`
	Balance       float64   `json:"balance"`
	UserId        string    `json:"user_id"`
	Flag          string    `json:"flag"`
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at"`
}

func (a *Account) ListAcount() (*Account, error) {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	if err := db.NewSelect().Model(a).Where("user_id = ? AND flag = ?", a.UserId, a.Flag).Scan(ctx); err != nil {
		return a, err
	}

	a.CreateAt = utils.LocalTime(a.CreateAt)
	a.UpdateAt = utils.LocalTime(a.UpdateAt)

	return a, nil
}

func (a *Account) ListAcounts() (*Account, error) {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	if err := db.NewSelect().Model(a).Where("user_id = ?", a.UserId).Scan(ctx); err != nil {
		return a, err
	}

	a.CreateAt = utils.LocalTime(a.CreateAt)
	a.UpdateAt = utils.LocalTime(a.UpdateAt)

	return a, nil
}
