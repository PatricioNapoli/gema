package models

import (
	"gema/server/utils"
	"github.com/go-pg/pg"
)

type Group struct {
	TableName struct{} `sql:"gema_groups"`

	BaseModel

	Name 	string
	Users   []*User `pg:",many2many:gema_membership"`
}

type Membership struct {
	TableName struct{} `sql:"gema_membership"`

	BaseModel

	UserId  int64 `sql:",pk"`
	User    *User
	GroupId int64 `sql:",pk"`
	Group   *Group
}

func FetchGroupByName(db *pg.DB, name string) *Group {
	group := &Group{}
	err := db.Model(group).Where("name = ?", name).Select()
	if err == pg.ErrNoRows {
		return nil
	}
	utils.Handle(err)

	return group
}