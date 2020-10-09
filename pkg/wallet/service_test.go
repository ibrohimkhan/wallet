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

	_, payments, _, err := s.addAcount(defaultTestAccount)
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

	_, _, _, err := s.addAcount(defaultTestAccount)
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

func TestService_FindFavoriteByID_success(t *testing.T) {
	s := newTestService()

	_, _, favorites, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favorite := favorites[0]
	result, err := s.FindFavoriteByID(favorite.ID)
	if err != nil {
		t.Errorf("FindFavoriteByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(favorite, result) {
		t.Errorf("FindFavoriteByID(): wrong favorite returned = %v", err)
		return
	}
}

func TestService_FindFavoriteByID_fail(t *testing.T) {
	s := newTestService()

	_, _, _, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindFavoriteByID(uuid.New().String())
	if err == nil {
		t.Error("FindFavoriteByID(): must return error, returned nil")
		return
	}

	if err != ErrFavoriteNotFound {
		t.Errorf("FindFavoriteByID(): must return ErrFavoriteNotFound, returned = %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()

	_, _, favorites, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favorite := favorites[0]
	_, err = s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): couldn't pay from favorites, error = %v", err)
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()

	_, payments, _, err := s.addAcount(defaultTestAccount)
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

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()

	_, payments, _, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	repay, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): can't repeat payment, error = %v", err)
		return
	}

	if repay.Category != payment.Category && repay.Amount != payment.Amount {
		t.Error("Repeat(): couldn't repeat the payment")
		return
	}
}

func TestService_Export_success(t *testing.T) {
	s := &Service{}

	account, err := s.RegisterAccount("+992937452945")
	if err != nil {
		t.Error(err)
		return
	}
	s.Deposit(account.ID, 100)

	account, err = s.RegisterAccount("+992937452946")
	if err != nil {
		t.Error(err)
		return
	}
	s.Deposit(account.ID, 101)

	account, err = s.RegisterAccount("+992937452947")
	if err != nil {
		t.Error(err)
		return
	}
	s.Deposit(account.ID, 102)

	err = s.ExportToFile("data/accounts.txt")
	if err != nil {
		t.Error(err)
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
	favorites	[]struct {
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
	favorites:		[]struct {
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

func (s *testService) addAcount(data testAccount) (*types.Account, []*types.Payment, []*types.Favorite, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	favorites := make([]*types.Favorite, len(data.favorites))
	for i := range data.favorites {
		paymentID := payments[i].ID
		favorites[i], err = s.FavoritePayment(paymentID, "favorite")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("can't make favorite, error = %v", err)
		}
	}

	return account, payments, favorites, nil
}