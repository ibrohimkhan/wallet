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