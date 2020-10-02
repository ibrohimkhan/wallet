package wallet

import (
	"github.com/google/uuid"
	"errors"
	"github.com/ibrohimkhan/wallet/v1.1.0/pkg/types"
)

// ErrPhoneRegistered - registration error
var ErrPhoneRegistered = errors.New("The given phone number already registered")

// ErrAmountMustBePositive - negative amount error
var ErrAmountMustBePositive = errors.New("The given amount must be greater than zero")

// ErrAccountNotFound - account does not exist
var ErrAccountNotFound = errors.New("Account not found")

// ErrNotEnoughBalance - balance is less then required for payment
var ErrNotEnoughBalance = errors.New("The balance does not have enough money")

// ErrPaymentNotFound - payment does not exist
var ErrPaymentNotFound = errors.New("Payment not found")

// ErrFavoriteNotFound - favorite does not exist
var ErrFavoriteNotFound = errors.New("Favorite not found")

// Service - storage for payments and accounts
type Service struct {
	nextAccountID	int64
	accounts 		[]*types.Account
	payments 		[]*types.Payment
	favorites		[]*types.Favorite
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

// Reject cencel payment
func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount

	return nil
}

// Repeat payment
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	repay, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}

	return repay, nil
}

// FavoritePayment creates favorite payment
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favoriteID := uuid.New().String()
	favorite := &types.Favorite {
		ID:			favoriteID,
		AccountID:	payment.AccountID,
		Name:		name,
		Amount:		payment.Amount,
		Category:	payment.Category,
	} 

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

// PayFromFavorite pay from favorite payment
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// FindAccountByID find account by id
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, ErrAccountNotFound
}

// FindPaymentByID searching payment by id
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

// FindFavoriteByID searching favorite by id
func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}