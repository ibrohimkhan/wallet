package wallet

import (
	"github.com/google/uuid"
	"fmt"
	"github.com/ibrohimkhan/wallet/v1.1.0/pkg/types"
	"reflect"
	"testing"
)

func TestService_FindAccountByID_success(t *testing.T) {
	s := newTestService()

	account, err := s.RegisterAccount("+992937452945")
	if err != nil {
		t.Errorf("FindAccountByID(): can't register new account, error = %v", err)
		return
	}

	result, err := s.FindAccountByID(account.ID)
	if err != nil {
		t.Errorf("FindAccountByID(): couldn't find account, error = %v", err)
		return
	}

	if !reflect.DeepEqual(account, result) {
		t.Error("FindAccountByID(): wrong account returned")
		return
	}
}

func TestService_FindAccountByID_accountNotFound(t *testing.T) {
	s := newTestService()
	_, err := s.FindAccountByID(1)

	if !reflect.DeepEqual(ErrAccountNotFound, err) {
		t.Errorf("FindAccountByID(): must return ErrAccountNotFound, but error = %v", err)
		return
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	result, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, result) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()

	_, _, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}

	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}

	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): account didn't changed, account = %v", savedAccount)
		return
	}
}

type testAccount struct {
	phone		types.Phone
	balance		types.Money
	payments	[]struct {
		amount		types.Money
		category	types.PaymentCategory
	}
}

var defaultTestAccount = testAccount {
	phone:			"+992937452945",
	balance:		10_000_00,
	payments:		[]struct {
		amount		types.Money
		category	types.PaymentCategory
	} {
		{ amount: 1_000_00, category:	"auto" },
	},
}

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{ Service: &Service{} }
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	return account, nil
}

func (s *testService) addAcount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}