package main

import (
	"fmt"
	"github.com/ibrohimkhan/wallet/v1.0.0/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992937452945")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(account.ID, 10)
	if err != nil {
		switch err  {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("The given amount must be positive value")
		case wallet.ErrAccountNotFound:
			fmt.Println("The user account does not exist")
		}

		return
	}
	
	fmt.Println(account.Balance)
}