package near

import (
	"bytes"
	"fmt"
	"knaq-wallet/config"
	"knaq-wallet/controller/erra"
	"os/exec"
	"strconv"
)

const yoctoNearQuantity = 1000000000000000000

var Cli info

type info struct {
	nearEnv              string
	nearID               string
	nearPrivateKey       string
	nearFtContract       string
	nearFtStorageDeposit string
}

func (i info) GenerateWallet(username string) (err error) {
	cmd := exec.Command("near", "generate-key", username)
	if err = cmd.Run(); err != nil {
		return erra.Error(err)
	}

	return nil
}

func (i info) RegisterToFtContract(publicKey string) (err error) {
	var (
		jsonString string
		cmd        *exec.Cmd
	)
	jsonString = fmt.Sprintf("{\"account_id\":\"%s\"}", publicKey)
	cmd = exec.Command("near", "call", config.Config.GetNearFtContract(), "storage_deposit", jsonString, "--account_id", config.Config.GetNearID(), "--amount", config.Config.GetNearFtStorageDeposit())
	if err = cmd.Run(); err != nil {
		return erra.Error(err)
	}

	return nil
}

func (i info) SendFt(out *bytes.Buffer, accountID string, amount string) (err error) {
	var (
		jsonString string
		cmd        *exec.Cmd
		amountF64  float64
		amountStr  string
	)
	amountF64, err = strconv.ParseFloat(amount, 64)
	if err != nil {
		return erra.Error(err)
	}
	amountF64 = amountF64 * yoctoNearQuantity
	amountStr = strconv.FormatFloat(amountF64, 'f', -1, 64)
	jsonString = fmt.Sprintf("{\"receiver_id\":\"%s\", \"amount\":\"%s\"}", accountID, amountStr)
	cmd = exec.Command("near", "call", config.Config.GetNearFtContract(), "ft_transfer", jsonString, "--accountId", config.Config.GetNearID(), "--depositYocto", "1")
	cmd.Stdout = out
	if err = cmd.Run(); err != nil {
		return erra.Error(err)
	}

	return nil
}

func init() {
	Cli = info{
		nearEnv:              config.Config.GetNearEnv(),
		nearID:               config.Config.GetNearID(),
		nearPrivateKey:       config.Config.GetNearPrivateKey(),
		nearFtContract:       config.Config.GetNearFtContract(),
		nearFtStorageDeposit: config.Config.GetNearFtStorageDeposit(),
	}
}
