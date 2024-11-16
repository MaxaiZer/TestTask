package dto

type WalletOperation struct {
	WalledID      string  `json:"walletId"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}
