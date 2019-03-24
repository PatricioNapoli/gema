package models

import (
	"gema/server/utils"
	"github.com/go-pg/pg"
)

type Token struct {
	TableName struct{} `sql:"gema_tokens"`

	BaseModel

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

func (t* Token) IsTokenValid(db *pg.DB) bool {
	err := db.Model(t).Where("token_hash = ?", t.TokenHash).Select()

	if err == pg.ErrNoRows {
		return false
	}

	return true
}
