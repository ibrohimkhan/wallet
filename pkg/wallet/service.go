package wallet

import (
	"path/filepath"
	"strings"
	"io"
	"strconv"
	"log"
	"os"
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

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return err
	}

	account.Balance += amount
	return nil
}

// Pay is a payment operation
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
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

	return s.Pay(payment.AccountID, payment.Amount, payment.Category)
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

	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
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

// ExportToFile saves accounts into a file
func (s *Service) ExportToFile(path string) error {
	err := s.exportAccountsToFile(path, "|")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ImportFromFile restores accounts into objects
func (s *Service) ImportFromFile(path string) error {
	data, err := s.getDataFromFile(path)
	if err != nil {
		log.Println(err)
		return err
	}

	accounts := s.parseStringToAccounts(data, "|")
	for _, account := range accounts {
		s.accounts = append(s.accounts, account)
	}

	return nil
}

// Export all available data (accounts, payments and favorites) to the given dir in files
func (s *Service) Export(dir string) error {
	if s.accounts != nil && len(s.accounts) > 0 {
		fullpath, err := s.getFullPath(dir, "accounts.dump")
		if err != nil {
			log.Println(err)
			return err
		}

		err = s.exportAccountsToFile(fullpath, "\n") 
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if s.payments != nil && len(s.payments) > 0 {
		fullpath, err := s.getFullPath(dir, "payments.dump")
		if err != nil {
			log.Println(err)
			return err
		}
		
		err = s.exportPaymentsToFile(fullpath)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if s.favorites != nil && len(s.favorites) > 0 {
		fullpath, err := s.getFullPath(dir, "favorites.dump")
		if err != nil {
			log.Println(err)
			return err
		}
		
		err = s.exportFavoritesToFile(fullpath)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

// Import all data from the given dir into objects such as accounts, payments and favorites
func (s *Service) Import(dir string) error {
	path, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	accountsPath := path + "/accounts.dump"
	if s.fileExist(accountsPath) {
		accounts, err := s.importAccountsFromFile(accountsPath)
		if err != nil {
			log.Println(err)
		}

		for _, dumpAccount := range accounts {
			if !s.containsAccount(dumpAccount, s.accounts) {
				s.accounts = append(s.accounts, dumpAccount)
				s.nextAccountID++
			}
		}

		log.Println("size of accounts = ", len(accounts))
	}

	paymentsPath := path + "/payments.dump"
	if s.fileExist(paymentsPath) {
		payments, err := s.importPaymentsFromFile(paymentsPath)
		if err != nil {
			log.Println(err)
		}

		for _, dumpPayment := range payments {
			if !s.containsPayment(dumpPayment, s.payments) {
				s.payments = append(s.payments, dumpPayment)
			}
		}

		log.Println("size of payments = ", len(payments))
	}

	favoritesPath := path + "/favorites.dump"
	if s.fileExist(favoritesPath) {
		favorites, err := s.importFavoritesFromFile(favoritesPath)
		if err != nil {
			log.Println(err)
		}

		for _, dumpFavorite := range favorites {
			if !s.containsFavorite(dumpFavorite, s.favorites) {
				s.favorites = append(s.favorites, dumpFavorite)
			}
		}

		log.Println("size of favorites = ", len(favorites))
	}

	return nil
}

// ExportAccountHistory get payments by accountid
func (s *Service) ExportAccountHistory(accountID int64) ([]*types.Payment, error) {
	var payments []*types.Payment

	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			payments = append(payments, payment)
		}
	}

	if payments == nil || len(payments) == 0 {
		return nil, ErrAccountNotFound
	}

	return payments, nil
}

// HistoryToFiles exports payments to files
func (s *Service) HistoryToFiles(payments []*types.Payment, dir string, records int) error {
	if payments == nil || len(payments) == 0 {
		return nil
	}

	if len(payments) <= records {
		fullpath, err := s.getFullPath(dir, "payments.dump")
		if err != nil {
			log.Println(err)
			return err
		}

		file, err := os.Create(fullpath)
		if err != nil {
			log.Println(err)
			return err
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.Println(err)
			}
		}()

		for _, payment := range payments {
			parsed := s.parsePaymentToString(payment)
			_, err := file.Write([]byte(parsed))
			if err != nil {
				log.Println(err)
				return err
			}
		}

	} else {
		count := 0

		filename := "payments" + strconv.Itoa(count + 1) + ".dump"
		fullpath, err := s.getFullPath(dir, filename)
		if err != nil {
			log.Println(err)
			return err
		}

		file, err := os.Create(fullpath)
		if err != nil {
			log.Println(err)
			return err
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.Println(err)
			}
		}()

		for index, payment := range payments {
			if index % records == 0 {
				count++

				filename = "payments" + strconv.Itoa(count) + ".dump"
				fullpath, err = s.getFullPath(dir, filename)
				if err != nil {
					log.Println(err)
					return err
				}

				file, err = os.Create(fullpath)
				if err != nil {
					log.Println(err)
					return err
				}
			}
			
			parsed := s.parsePaymentToString(payment)
			_, err = file.Write([]byte(parsed))
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}

// GetPayments returns payments
func (s *Service) GetPayments() []*types.Payment {
	return s.payments
}

func (s *Service) fileExist(path string) bool {
	info, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
	}
	
    return !info.IsDir()
}

func (s *Service) getFullPath(dir string, filename string) (string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	} 
	
	fullpath := path + "/" + filename
	return fullpath, nil
}

func (s *Service) exportAccountsToFile(path string, sep string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	for _, account := range s.accounts {
		parsed := s.parseAccountToString(account, sep)
		_, err := file.Write([]byte(parsed))
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (s *Service) exportPaymentsToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	for _, payment := range s.payments {
		parsed := s.parsePaymentToString(payment)
		_, err := file.Write([]byte(parsed))
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (s *Service) exportFavoritesToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	for _, favorite := range s.favorites {
		parsed := s.parseFavoriteToString(favorite)
		_, err := file.Write([]byte(parsed))
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (s *Service) importAccountsFromFile(path string) ([]*types.Account, error) {
	data, err := s.getDataFromFile(path)
	if err != nil {
		return nil, err
	}

	accounts := s.parseStringToAccounts(data, "\n")
	return accounts, nil
}

func (s *Service) importPaymentsFromFile(path string) ([]*types.Payment, error) {
	data, err := s.getDataFromFile(path)
	if err != nil {
		return nil, err
	}

	payments := s.parseStringToPayments(data)
	return payments, nil
}

func (s *Service) importFavoritesFromFile(path string) ([]*types.Favorite, error) {
	data, err := s.getDataFromFile(path)
	if err != nil {
		return nil, err
	}

	favorites := s.parseStringToFavorites(data)
	return favorites, nil
}

func (s *Service) getDataFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4096)

	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			return "", err
		}

		content = append(content, buf[:read]...)
	}

	data := string(content)
	return data, nil
}

func (s *Service) parseAccountToString(account *types.Account, sep string) string {
	parsed := strconv.FormatInt(account.ID, 10) + ";"
	parsed += string(account.Phone) + ";"
	parsed += strconv.FormatInt(int64(account.Balance), 10) + sep
	
	return parsed
}

func (s *Service) parseStringToAccounts(data string, sep string) []*types.Account {
	var accounts []*types.Account

	data = strings.TrimSpace(data)
	for _, items := range strings.Split(data, sep) {
		item := strings.Split(items, ";")

		if len(item) < 3 {
			account := &types.Account {
				ID:			0,
				Phone:		types.Phone(item[0]),
				Balance:	0,
			}

			accounts = append(accounts, account)
		
		} else {
			id, _ 		:= strconv.ParseInt(item[0], 10, 64)
			phone 		:= types.Phone(item[1])
			balance, _ 	:= strconv.ParseInt(item[2], 10, 64)

			account := &types.Account {
				ID:			id,
				Phone:		phone,
				Balance:	types.Money(balance),
			}

			accounts = append(accounts, account)
		}
	}

	return accounts
}

func (s *Service) parsePaymentToString(payment *types.Payment) string {
	parsed := payment.ID + ";"
	parsed += strconv.FormatInt(payment.AccountID, 10) + ";"
	parsed += strconv.FormatInt(int64(payment.Amount), 10) + ";"
	parsed += string(payment.Category) + ";"
	parsed += string(payment.Status) + "\n"

	return parsed
}

func (s *Service) parseStringToPayments(data string) []*types.Payment {
	var payments []*types.Payment

	data = strings.TrimSpace(data)
	for _, items := range strings.Split(data, "\n") {
		item := strings.Split(items, ";")

		id			 	:= string(item[0])
		accountID, _ 	:= strconv.ParseInt(item[1], 10, 64)
		amount, _ 		:= strconv.ParseInt(item[2], 10, 64)
		category		:= string(item[3])
		status			:= string(item[4])

		payment := &types.Payment {
			ID:				id,
			AccountID:		accountID,
			Amount:			types.Money(amount),
			Category:		types.PaymentCategory(category),
			Status:			types.PaymentStatus(status),
		}

		payments = append(payments, payment)
	}

	return payments
}

func (s *Service) parseFavoriteToString(favorite *types.Favorite) string {
	parsed := favorite.ID + ";"
	parsed += strconv.FormatInt(favorite.AccountID, 10) + ";"
	parsed += favorite.Name + ";"
	parsed += strconv.FormatInt(int64(favorite.Amount), 10) + ";"
	parsed += string(favorite.Category) + "\n"

	return parsed
}

func (s *Service) parseStringToFavorites(data string) []*types.Favorite {
	var favorites []*types.Favorite

	data = strings.TrimSpace(data)
	for _, items := range strings.Split(data, "\n") {
		item := strings.Split(items, ";")

		id			 	:= string(item[0])
		accountID, _ 	:= strconv.ParseInt(item[1], 10, 64)
		name			:= string(item[2])
		amount, _ 		:= strconv.ParseInt(item[3], 10, 64)
		category		:= string(item[4])

		favorite := &types.Favorite {
			ID:				id,
			AccountID:		accountID,
			Name:			name,
			Amount:			types.Money(amount),
			Category:		types.PaymentCategory(category),
		}

		favorites = append(favorites, favorite)
	}

	return favorites
}

func (s *Service) containsAccount(item *types.Account, items []*types.Account) bool {
	for _, value := range items {
		if value.ID == item.ID {
			value.Phone 	= item.Phone
			value.Balance 	= item.Balance

			return true
		}
	}

	return false
}

func (s *Service) containsPayment(item *types.Payment, items []*types.Payment) bool {
	for _, value := range items {
		if value.ID == item.ID {
			value.AccountID 	= item.AccountID
			value.Amount 		= item.Amount
			value.Category	 	= item.Category
			value.Status 		= item.Status

			return true
		}
	}

	return false
}

func (s *Service) containsFavorite(item *types.Favorite, items []*types.Favorite) bool {
	for _, value := range items {
		if value.ID == item.ID {
			value.AccountID 	= item.AccountID
			value.Name 			= item.Name
			value.Amount 		= item.Amount
			value.Category	 	= item.Category

			return true
		}
	}

	return false
}