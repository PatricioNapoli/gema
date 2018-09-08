package database

import (
	"gema/server/models"

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

	err := d.migrate()
	if err != nil {
		panic(err)
	}

	return d
}

func (s *Database) Dispose() {
	s.SQL.Close()
}

func (s *Database) IsFirstLogin() bool {
	count, err := s.SQL.Model((*models.User)(nil)).Count()
	if err != nil {
		panic(err)
	}
	return count == 0
}

func (s *Database) migrate() error {
	for _, model := range []interface{}{(*models.User)(nil)} {
		err := s.SQL.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
