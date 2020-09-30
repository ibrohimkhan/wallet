package wallet

import (
	"github.com/google/uuid"
	"errors"
	"github.com/ibrohimkhan/wallet/pkg/types"
)

// ErrPhoneRegistered - registration error
var ErrPhoneRegistered = errors.New("The given phone number already registered")

// ErrAmountMustBePositive - negative amount error
var ErrAmountMustBePositive = errors.New("The given amount must be greater than zero")

// ErrAccountNotFound - account does not exist
var ErrAccountNotFound = errors.New("Account not found")

// ErrNotEnoughBalance - balance is less then required for payment
var ErrNotEnoughBalance = errors.New("The balance does not have enough money")

// Service - storage for payments and accounts
type Service struct {
	nextAccountID	int64
	accounts 		[]*types.Account
	payments 		[]*types.Payment
}

// RegisterAccount registering new account
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account {
		ID:			s.nextAccountID,
		Phone:		phone,
		Balance:	0,
	}

	s.accounts = append(s.accounts, account)
	return account, nil
}

// Deposit add money based on account id
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, item := range s.accounts {
		if item.ID == accountID {
			account = item
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

// Pay is a payment operation
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, item := range s.accounts {
		if item.ID == accountID {
			account = item
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()

	payment := &types.Payment {
		ID:			paymentID,
		AccountID:	accountID,
		Amount:		amount,
		Category:	category,
		Status:		types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}

// FindAccountByID find account by id
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account
	for _, item := range s.accounts {
		if item.ID == accountID {
			account = item
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}