package models

type User struct {
	TableName struct{} `sql:"gema_users"`

	Id int64
	Email string
	Name string
	Hash string
}

type User2 struct {
	TableName struct{} `sql:"gema_users"`

	Id int64
	Email string
	Name string
	Hash string
}