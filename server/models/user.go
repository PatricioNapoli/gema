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

func (s *User) FetchUserByEmail(db *pg.DB, email string) {
	err := db.Model(s).Where("email = ?", email).Select()
	if err != nil {
		panic(err)
	}
}
