package models

import (
	"github.com/go-pg/pg"
)

type User struct {
	TableName struct{} `sql:"gema_users"`

	Id    int64
	Email string
	Name  string
	Hash  string
}

func FetchUserByEmail(db *pg.DB, email string) *User {
	user := &User{}
	err := db.Model(user).Where("email = ?", email).Select()
	if err == pg.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}

	return user
}
