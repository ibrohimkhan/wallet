package wallet

import "github.com/ibrohimkhan/wallet/pkg/types"

// Service - storage for payments and accounts
type Service struct {
	nextAccountID	int64
	accounts 		[]types.Account
	payments 		[]types.Payment
}