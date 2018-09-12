package models

import (
	"github.com/go-pg/pg"
	"gema/server/utils"
)

type User struct {
	TableName struct{} `sql:"gema_users"`

	Id    int64
	Email string
	Name  string
	Hash  string
}

// Fetches the user based on an email filter, returns nil if not found.
func FetchUserByEmail(db *pg.DB, email string) *User {
	user := &User{}
	err := db.Model(user).Where("email = ?", email).Select()
	if err == pg.ErrNoRows {
		return nil
	}
	utils.Handle(err)

	return user
}
