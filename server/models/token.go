package models

import (
	"gema/server/security"
	"gema/server/utils"
	"github.com/go-pg/pg"
)

type Token struct {
	TableName struct{} `sql:"gema_tokens"`

	BaseModel

	User string
	TokenHash string
}

func (t *Token) InsertToken(db *pg.DB) int64 {
	var uid int64
	_, err := db.Model(t).Returning("id", uid).Insert()

	if err != nil {
		utils.Handle(err)
	}

	return uid
}

func (t* Token) VerifyUserPassword(db *pg.DB, password string) bool {
	err := db.Model(t).Select()

	if err == pg.ErrNoRows {
		return false
	}

	return security.VerifyPassword(t.TokenHash, password)
}
