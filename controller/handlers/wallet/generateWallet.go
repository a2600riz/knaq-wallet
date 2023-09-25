package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"knaq-wallet/config"
	"knaq-wallet/controller/erra"
	"knaq-wallet/near"
	"os"
	"os/exec"
	"strings"
)

func GenerateWallet(username string) (nearWallet NearWalletGeneration, err error) {
	var (
		cmd            *exec.Cmd
		homeDir        string
		walletJsonFile []byte
		filePath       string
		pbcKeys        []string
		b58decPubKey   []byte
	)
	if err = near.Cli.GenerateWallet(username); err != nil {
		return nearWallet, erra.Error(err)
	}
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return nearWallet, erra.Error(err)
	}
	filePath = fmt.Sprintf("%s/.near-credentials/%s/%s.json", homeDir, config.Config.GetNearEnv(), username)

	walletJsonFile, err = os.ReadFile(filePath)
	if err != nil {
		return nearWallet, erra.Error(err)
	}
	if err = json.Unmarshal(walletJsonFile, &nearWallet); err != nil {
		return nearWallet, erra.Error(err)
	}

	pbcKeys = strings.Split(nearWallet.PublicKey, ":")
	b58decPubKey = base58.Decode(pbcKeys[1])
	nearWallet.PublicKey = hex.EncodeToString(b58decPubKey)
	if err = near.Cli.RegisterToFtContract(nearWallet.PublicKey); err != nil {
		return nearWallet, erra.Error(err)
	}

	cmd = exec.Command("rm", "-rf", filePath)
	if err = cmd.Run(); err != nil {
		return nearWallet, erra.Error(err)
	}

	return nearWallet, nil
}
