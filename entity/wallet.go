package entity

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
	"knaq-wallet/tools/security"
)

const (
	SaltBytes = 32
	Iteration = 4096
)

type Wallet struct {
	Common
	UserID     uint64 `gorm:"unique;not null"`
	AccountID  string `gorm:"type:varchar(128);not null"`
	PublicKey  string `gorm:"type:varchar(128);not null"`
	PrivateKey string `gorm:"type:text;not null"`
}

func (b *Wallet) TableName() string {
	return "wallets"
}

func (b *Wallet) GetTx() *gorm.DB {
	return b.tx
}

func (b *Wallet) InitTx() error {
	err := b.tx.Error
	b.tx = nil
	return err
}

func (b *Wallet) SetTx(tx *gorm.DB) *Wallet {
	b.tx = tx
	return b
}

func (b *Wallet) SelectAll() (result []Wallet, err error) {
	if err = db.Find(&result).Error; err != nil {
		return
	}

	return
}

func (b *Wallet) CheckEncrypted() (bool, error) {
	if err := db.Where("user_id=?", b.UserID).Where("deleted_at is null").First(b).Error; err != nil {
		return false, err
	}
	if b.ID == 0 {
		return false, nil
	}
	return true, nil
}

func (b *Wallet) SelectOneByUserId(email, password string) error {
	if err := db.Where("user_id=?", b.UserID).First(b).Error; err != nil {
		return err
	}

	hashedKey := pbkdf2.Key([]byte(password), []byte(email), Iteration, SaltBytes, sha256.New)

	enc := security.Encryption{SecretKey: hashedKey}

	pk := enc.Decrypt(b.PrivateKey)
	if pk == "" {
		return fmt.Errorf("couldn't decrypt private key")
	}

	return nil
}

func (b *Wallet) SelectKeysByUserId() error {

	if err := db.Where("user_id=?", b.UserID).First(b).Error; err != nil {
		return err
	}

	return nil
}

func (b *Wallet) Create(email, password string) *Wallet {
	if b.tx == nil {
		b.tx = db.Begin()
	}

	hashedKey := pbkdf2.Key([]byte(password), []byte(email), Iteration, SaltBytes, sha256.New)
	enc := security.Encryption{SecretKey: hashedKey}
	encryptedString, err := enc.Encrypt(b.PrivateKey)
	if err != nil {
		b.tx.Error = err
		return b
	}

	b.PrivateKey = encryptedString

	if err = b.tx.Create(b).Error; err != nil {
		b.tx.Error = err
	}

	return b
}

func (b *Wallet) UpdateById(email, password string) *Wallet {
	if b.tx == nil {
		b.tx = db.Begin()
	}

	hashedKey := pbkdf2.Key([]byte(password), []byte(email), Iteration, SaltBytes, sha256.New)
	enc := security.Encryption{SecretKey: hashedKey}
	encryptedString, err := enc.Encrypt(b.PrivateKey)
	if err != nil {
		b.tx.Error = err
		return b
	}

	b.PrivateKey = encryptedString

	if err = b.tx.Updates(b).Error; err != nil {
		b.tx.Error = err
	}

	return b
}

func (b *Wallet) Commit() error {
	if b.tx.Error != nil {
		b.tx.Rollback()
		return b.InitTx()
	}
	if b.tx == nil {
		return fmt.Errorf("db transaction is nil")
	}
	b.tx.Error = b.tx.Commit().Error
	return b.InitTx()
}
