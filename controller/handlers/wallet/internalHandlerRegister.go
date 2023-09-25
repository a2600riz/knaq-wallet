package wallet

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"knaq-wallet/controller/erra"
	"knaq-wallet/entity"
	"knaq-wallet/near"
	"net/http"
	"strings"
)

type InternalWalletHandler struct{}

func (_ InternalWalletHandler) Register(group *echo.Group) {
	walletGroup := group.Group(baseGroupPath)
	walletGroup.POST("/knaq/ft/send", sendFromKnaq)
}

func sendFromKnaq(ctx echo.Context) error {
	var (
		err                 error
		request             SendTokenRequest
		out                 = bytes.Buffer{}
		txID                string
		ftTransactionEntity entity.FtTransaction
		resp                map[string]interface{}
		respStatus          int
	)

	if err = ctx.Bind(&request); err != nil {
		return erra.BadRequest(err)
	}
	if err = near.Cli.SendFt(&out, request.AccountID, request.Amount); err != nil {
		return erra.Error(err)
	}
	{
		idx1 := strings.Index(out.String(), "Transaction Id ")
		idx2 := strings.Index(out.String(), "To see the transaction")
		if idx1 >= 0 && idx2-1 <= len(out.String()) {
			str := out.String()[idx1 : idx2-1]
			txID = strings.TrimPrefix(str, "Transaction Id ")

			ftTransactionEntity = entity.FtTransaction{
				UserID:        request.UserID,
				ReceiverID:    request.AccountID,
				TransactionID: txID,
				Amount:        request.Amount,
			}

			resp = map[string]interface{}{
				"tx_id": txID,
			}

			respStatus = http.StatusCreated
		} else {
			fmt.Println(out.String())

			ftTransactionEntity = entity.FtTransaction{
				UserID:     request.UserID,
				ReceiverID: request.AccountID,
				Amount:     request.Amount,
				Error: sql.NullString{
					String: out.String(),
					Valid:  true,
				},
			}

			resp = map[string]interface{}{
				"error_message": "couldn't proceed the transaction. please check address and amount if it's correct.",
			}

			respStatus = http.StatusBadRequest
		}
	}

	if err = ftTransactionEntity.Create().Commit(); err != nil {
		return erra.Error(err)
	}

	return ctx.JSON(respStatus, resp)
}
