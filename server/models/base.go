package models

import (
	"context"
	"github.com/go-pg/pg/orm"
	"time"
)

type BaseModel struct {
	Id int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *BaseModel) BeforeInsert(c context.Context, db orm.DB) error {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}
	return nil
}

func (m *BaseModel) BeforeUpdate(c context.Context, db orm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
