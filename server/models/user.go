package models

import "github.com/go-pg/pg"

type User struct {
	TableName struct{} `sql:"gema_users"`

	Id int64
	Email string
	Name string
	Hash string
}

func GetUser(db *pg.DB, email string) *User{
	u := &User{}
	err := db.Model(u).Where("email = ?", email).Select()
	if err != nil {
		return nil
	}
	return u
}