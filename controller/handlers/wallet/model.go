package wallet

type CreateWalletRequest struct {
	Password string `json:"password"`
}

type SendTokenRequest struct {
	UserID    uint64 `json:"user_id"`
	AccountID string `json:"account_id"`
	Amount    string `json:"amount"`
}

type SendTokenResponse struct {
	TxID   string `json:"tx_id"`
	Status string `json:"status"`
}
