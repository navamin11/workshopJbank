package models

import (
	"time"
	"workshop/configs"
	"workshop/utils"

	"github.com/twinj/uuid"
	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`
	Id            string    `json:"id"`
	AccountNo     string    `json:"account_no"`
	AccountName   string    `json:"account_name"`
	Bank          string    `json:"bank"`
	Type          string    `json:"type"`
	MoneyIn       float64   `json:"money_in"`
	MoneyOut      float64   `json:"money_out"`
	DepositId     string    `json:"deposit_id"`
	UpdateBy      string    `json:"update_by"`
	UpdateAt      time.Time `json:"update_at"`
}

func (trans *Transaction) TransactionIn() error {
	var account Account

	account.Id = trans.DepositId

	if err := account.addMoney(trans.MoneyIn); err != nil {
		return err
	}

	if err := trans.insertTransactionIn(); err != nil {
		return err
	}

	return nil
}

func (trans *Transaction) TransactionOut() error {
	var account Account

	account.Id = trans.DepositId

	if err := account.subMoney(trans.MoneyOut); err != nil {
		return err
	}

	if err := trans.insertTransactionOut(); err != nil {
		return err
	}

	return nil
}

func (trans *Transaction) insertTransactionIn() error {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	trans.Id = uuid.NewV4().String()
	trans.UpdateBy = "thirdparty"
	trans.UpdateAt = utils.LocalTime(time.Now())

	if _, err := db.NewInsert().Model(trans).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (a *Account) addMoney(amount float64) error {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	var balance float64
	if err := db.NewSelect().Model(a).Where("id = ? ", a.Id).Scan(ctx); err != nil {
		return err
	}

	balance = a.Balance + amount
	a.Balance = balance
	a.UpdateAt = utils.LocalTime(time.Now())

	if _, err := db.NewUpdate().Model(a).Column("balance", "update_at").Where("id = ? ", a.Id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (trans *Transaction) insertTransactionOut() error {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	trans.Id = uuid.NewV4().String()
	trans.UpdateBy = "owner"
	trans.UpdateAt = utils.LocalTime(time.Now())

	if _, err := db.NewInsert().Model(trans).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (a *Account) subMoney(amount float64) error {
	db := configs.ConnectPostgreSQL()
	defer db.Close()

	var balance float64

	if err := db.NewSelect().Model(a).Where("id = ? AND balance > ? AND balance >= ?", a.Id, 0, amount).Scan(ctx); err != nil {
		return err
	}

	balance = a.Balance - amount
	a.Balance = balance
	a.UpdateAt = utils.LocalTime(time.Now())

	if _, err := db.NewUpdate().Model(a).Column("balance", "update_at").Where("id = ? AND balance > ? ", a.Id, 0).Exec(ctx); err != nil {
		return err
	}

	return nil
}
