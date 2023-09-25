package wallet

import "encoding/json"

type NearWalletGeneration struct {
	AccountId  string `json:"account_id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func (n *NearWalletGeneration) Unmarshal(data []byte) error {
	return json.Unmarshal(data, n)
}
