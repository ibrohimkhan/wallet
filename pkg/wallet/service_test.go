package wallet

import (
	"reflect"
	"testing"
)

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+992937452945")
	
	if err != nil {
		t.Error(err)
	}

	result, err := svc.FindAccountByID(account.ID)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(account, result) {
		t.Error("Accounts are different!")
	}
}

func TestService_FindAccountByID_accountNotFound(t *testing.T) {
	svc := &Service{}
	_, err := svc.FindAccountByID(1)

	if !reflect.DeepEqual(ErrAccountNotFound, err) {
		t.Error("Invalid result")
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+992937452945")
	err = svc.Deposit(account.ID, 100)
	
	if err != nil {
		t.Error(err)
	}

	payment, err := svc.Pay(account.ID, 15, "auto")
	if err != nil {
		t.Error(err)
	}

	result, err := svc.FindPaymentByID(payment.ID)

	if !reflect.DeepEqual(payment, result) {
		t.Error("Invalid result")
	}
}

func TestService_FindPaymentByID_notFound(t *testing.T) {
	svc := &Service{}
	_, err := svc.FindPaymentByID("110")

	if !reflect.DeepEqual(ErrPaymentNotFound, err) {
		t.Error("Invalid result")
	}
}

func TestService_Reject_success(t *testing.T) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+992937452945")
	err = svc.Deposit(account.ID, 100)
	
	if err != nil {
		t.Error(err)
	}

	payment, err := svc.Pay(account.ID, 15, "auto")
	if err != nil {
		t.Error(err)
	}

	err = svc.Reject(payment.ID)
	if err != nil {
		t.Error(err)
	}

	if account.Balance != 100 {
		t.Error("invalid result")
	}
}