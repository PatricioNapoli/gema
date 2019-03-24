package database

import (
	"gema/server/models"

	"gema/server/utils"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Database struct {
	SQL *pg.DB
}

func New(sql *pg.DB) *Database {
	d := &Database{
		SQL: sql,
	}

	return d
}

func (s *Database) Dispose() {
	if err := s.SQL.Close(); err != nil {
		utils.Handle(err)
	}
}

func (s *Database) IsFirstUser() bool {
	count, err := s.SQL.Model((*models.User)(nil)).Count()
	utils.Handle(err)
	return count == 0
}

func (s *Database) IsFirstBoot() bool {
	count, err := s.SQL.Model((*models.Group)(nil)).Count()
	utils.Handle(err)
	return count == 0
}

func (s *Database) loadInitialData() {
	groups := []string{"master", "dev", "client"}

	for _, group := range groups {
		_, err := s.SQL.Model(&models.Group{
			Name: group,
		}).Insert()

		if err != nil {
			utils.Handle(err)
		}
	}
}

func (s *Database) Migrate() bool {
	tables := []interface{}{
		(*models.User)(nil),
		(*models.Group)(nil),
		(*models.Membership)(nil),
		(*models.Token)(nil),
	}

	for _, model := range tables {
		err := s.SQL.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		utils.Handle(err)
	}

	isFirst := s.IsFirstBoot()

	if isFirst {
		s.loadInitialData()
	}

	return isFirst
}
