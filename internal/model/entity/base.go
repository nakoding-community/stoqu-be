package entity

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/ctxval"
	"gorm.io/gorm"
)

type Entity struct {
	ID string `json:"id" gorm:"primaryKey;"`

	CreatedAt  *time.Time `json:"created_at"`
	CreatedBy  string     `json:"created_by"`
	ModifiedAt *time.Time `json:"modified_at"`
	ModifiedBy string     `json:"modified_by,omitempty"`
}

func (m *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	now := time.Now()
	if m.CreatedAt == nil {
		m.CreatedAt = &now
	}
	m.CreatedBy = constant.DB_DEFAULT_SYSTEM
	if m.ModifiedAt == nil {
		m.ModifiedAt = &now
	}
	m.ModifiedBy = constant.DB_DEFAULT_SYSTEM
	return
}

func (m *Entity) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	if m.ModifiedAt == nil {
		m.ModifiedAt = &now
	}
	m.ModifiedBy = constant.DB_DEFAULT_SYSTEM

	authCtx := ctxval.GetAuthValue(tx.Statement.Context)
	if authCtx != nil {
		m.ModifiedBy = authCtx.Name
	}
	return
}

type MasterEntity struct {
	Code string `json:"code" gorm:"not null;unique;size:50"`
	Name string `json:"name" gorm:"not null;unique;size:50"`
}

type Tabler interface {
	TableName() string
}
