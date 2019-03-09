package models

import (
	"gema/server/utils"
	"github.com/go-pg/pg"
	"time"
)

type User struct {
	TableName struct{} `sql:"gema_users"`

	BaseModel

	Email string
	Name  string
	Hash  string

	LastSignIn time.Time

	Groups   []*Group `pg:",many2many:gema_membership"`
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
