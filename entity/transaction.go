package entity

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

type FtTransaction struct {
	Common
	UserID        uint64         `gorm:"not null"`
	ReceiverID    string         `gorm:"type:varchar(128);not null"`
	TransactionID string         `gorm:"type:varchar(128)"`
	Amount        string         `gorm:"type:varchar(255);not null"`
	Error         sql.NullString `gorm:"type:text"`
}

func (t *FtTransaction) TableName() string {
	return "ft_transactions"
}

func (t *FtTransaction) GetTx() *gorm.DB {
	return t.tx
}

func (t *FtTransaction) InitTx() error {
	err := t.tx.Error
	t.tx = nil
	return err
}

func (t *FtTransaction) SetTx(tx *gorm.DB) *FtTransaction {
	t.tx = tx
	return t
}

func (t *FtTransaction) SelectAll() (result []FtTransaction, err error) {
	if err = db.Find(&result).Error; err != nil {
		return
	}

	return
}

func (t *FtTransaction) SelectOneByUserId() error {
	if err := db.Where("user_id=?", t.UserID).First(t).Error; err != nil {
		return err
	}

	return nil
}

func (t *FtTransaction) Create() *FtTransaction {
	if t.tx == nil {
		t.tx = db.Begin()
	}

	if err := t.tx.Create(t).Error; err != nil {
		t.tx.Error = err
	}

	return t
}

func (t *FtTransaction) UpdateById() *FtTransaction {
	if t.tx == nil {
		t.tx = db.Begin()
	}

	if err := t.tx.Updates(t).Error; err != nil {
		t.tx.Error = err
	}

	return t
}

func (t *FtTransaction) Commit() error {
	if t.tx.Error != nil {
		t.tx.Rollback()
		return t.InitTx()
	}
	if t.tx == nil {
		return fmt.Errorf("db transaction is nil")
	}
	t.tx.Error = t.tx.Commit().Error
	return t.InitTx()
}
