package wallet

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"knaq-wallet/controller/erra"
	"knaq-wallet/entity"
	"knaq-wallet/tools/security"
	"net/http"
)

type SecureWalletHandler struct{}

const baseGroupPath = "/near"

func (_ SecureWalletHandler) Register(group *echo.Group) {
	walletGroup := group.Group(baseGroupPath)
	walletGroup.GET("/key/check", checkNearKey)
	walletGroup.GET("/key", getNearKey)
	walletGroup.POST("/key", createNearKey)
}

func checkNearKey(ctx echo.Context) error {
	var (
		err          error
		claimData    = ctx.Get(security.ClaimDataKey).(security.UserClaimData)
		walletEntity = entity.Wallet{UserID: claimData.UserID}
		hasEncrypted bool
	)

	hasEncrypted, err = walletEntity.CheckEncrypted()
	if err != nil {
		return erra.Error(err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"has_encrypted": hasEncrypted,
	})
}
func getNearKey(ctx echo.Context) error {
	var (
		err                     error
		claimData               = ctx.Get(security.ClaimDataKey).(security.UserClaimData)
		walletEncryptionRequest CreateWalletRequest
		walletEntity            entity.Wallet
	)

	if err = ctx.Bind(&walletEncryptionRequest); err != nil {
		return erra.BadRequest(err)
	}

	walletEntity = entity.Wallet{
		UserID: claimData.UserID,
	}
	if err = walletEntity.SelectKeysByUserId(); err != nil {
		return erra.Error(err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"encoded_public_key":    walletEntity.PublicKey,
		"encrypted_private_key": walletEntity.PrivateKey,
	})
}
func createNearKey(ctx echo.Context) error {
	var (
		err                     error
		claimData               = ctx.Get(security.ClaimDataKey).(security.UserClaimData)
		hasEncrypted            bool
		username                = uuid.NewString()
		walletEncryptionRequest CreateWalletRequest
		nearWallet              NearWalletGeneration
		walletEntity            = entity.Wallet{UserID: claimData.UserID}
	)

	if err = ctx.Bind(&walletEncryptionRequest); err != nil {
		return erra.BadRequest(err)
	}

	hasEncrypted, err = walletEntity.CheckEncrypted()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return erra.Error(err)
		}
	}

	if hasEncrypted {
		return erra.BadRequest(fmt.Errorf("key exists"))
	}

	nearWallet, err = GenerateWallet(username)
	if err != nil {
		return erra.Error(err)
	}

	walletEntity.AccountID = nearWallet.AccountId
	walletEntity.PublicKey = nearWallet.PublicKey
	walletEntity.PrivateKey = nearWallet.PrivateKey
	if err = walletEntity.Create(claimData.Email, walletEncryptionRequest.Password).Commit(); err != nil {
		return erra.Error(err)
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"encoded_public_key":    nearWallet.PublicKey,
		"encrypted_private_key": walletEntity.PrivateKey,
	})
}
