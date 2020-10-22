package wallet

import (
	"strings"
	"os"
	"path/filepath"
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

func  TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()

	_, payments, _, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	fav, err := s.FavoritePayment(payment.ID, string(payment.Category))
	if err  != nil {
		t.Error(err)
		return
	}

	if fav.Amount != payment.Amount {
		t.Error(err)
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

	err = s.ExportToFile("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}

	err = os.Remove("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Import_success(t *testing.T) {
	s := &Service{}

	if len(s.accounts) > 0 {
		t.Fail()
		return
	}

	account, err := s.RegisterAccount("+992937452947")
	if err != nil {
		t.Error(err)
		return
	}
	s.Deposit(account.ID, 102)

	err = s.ExportToFile("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}

	err = s.ImportFromFile("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}

	if len(s.accounts) == 0 {
		t.Fail()
		return
	}

	err = os.Remove("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ConcurrentSumOfPayments_success(t *testing.T) {
	s := newTestService()

	payments := []*types.Payment {
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
		{ Amount: 1, Category:	"auto" },
	}
	
	s.payments = payments

	expected := types.Money(15)
	result := s.SumPayments(3)

	if result != expected {
		t.Errorf("invalid result! Expected %v, got %v", expected, result)
		return
	}
}

func TestService_ConcurrentSumOfPayments_fail(t *testing.T) {
	s := newTestService()

	payments := []*types.Payment {
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 4, Amount: 1, Category: "auto" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 3, Amount: 1, Category: "auto" },
	}
	
	s.payments = payments

	expected := types.Money(15)
	result := s.SumPayments(3)

	if result == expected {
		t.Errorf("invalid result! Expected %v, got %v", expected, result)
		return
	}
}

func TestService_FilterPayments_success(t *testing.T) {
	s := newTestService()

	accounts := []*types.Account {
		{ ID: 1, Phone: "111111", Balance: 0 },
	}

	s.accounts = accounts

	payments := []*types.Payment {
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
	}
	
	s.payments = payments
	filtered, err := s.FilterPayments(1, 2)
	if err != nil {
		t.Error(err)
		return
	}

	if len(filtered) != 10 {
		t.Fail()
		t.Errorf("invalid result! Expected %v, got %v", 8, len(filtered))
		return
	}
}

func TestService_FilterPayments_withNoPayments(t *testing.T) {
	s := newTestService()

	accounts := []*types.Account {
		{ ID: 1, Phone: "111111", Balance: 0 },
	}

	s.accounts = accounts

	filtered, err := s.FilterPayments(1, 1)
	if err != nil {
		t.Error(err)
		return
	}

	if filtered != nil {
		t.Fail()
	}
}

func TestService_FilterPaymentsByFn_success(t *testing.T) {
	s := newTestService()

	payments := []*types.Payment {
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "book" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "book" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 4, Amount: 1, Category: "auto" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 3, Amount: 1, Category: "book" },
	}
	
	s.payments = payments
	filter := s.filter
	filtered, err := s.FilterPaymentsByFn(filter, 3)
	if err != nil {
		t.Error(err)
		return
	}

	want := 3
	if len(filtered) != want {
		t.Fail()
		t.Errorf("invalid result! Expected %v, got %v", want, len(filtered))
		return
	}
}

func TestService_SumPaymentsWithProgress_success(t *testing.T) {
	s := newTestService()

	var payments []*types.Payment
	for i := 0; i < 100000000; i++ {
		payment := &types.Payment {
			AccountID: 1, Amount: 1, Category: "auto",
		}

		payments = append(payments, payment)
	}

	s.payments = payments
	ch := s.SumPaymentsWithProgress()

	total := types.Money(0)
	for i := 0; i < 100000000; i++ {
		progress := <- ch
		total += progress.Result
	}

	if total != 100000000 {
		t.Errorf("invalid result! Expected %v, got %v", 100000000, total)
	}
}

func TestService_SumPaymentsWithProgress_empty(t *testing.T) {
	s := newTestService()
	ch := s.SumPaymentsWithProgress()
	if 0 != len(ch) {
		t.Fail()
	}
}

func TestService_SumOf_success(t *testing.T) {
	s := newTestService()

	payments := []*types.Payment {
		{ Amount: 1_000_00, Category:	"auto" },
		{ Amount: 1_000_00, Category:	"auto" },
		{ Amount: 1_000_00, Category:	"auto" },
	}

	expected := types.Money(3_000_00)
	result := s.sumOf(payments)

	if result != expected {
		t.Error("invalid result")
		return
	}
}

func TestService_SumOf_fail(t *testing.T) {
	s := newTestService()

	payments := []*types.Payment {
		{ Amount: 1_000_00, Category:	"auto" },
		{ Amount: 1_000_00, Category:	"auto" },
	}

	expected := types.Money(3_000_00)
	result := s.sumOf(payments)

	if result == expected {
		t.Error("invalid result")
		return
	}
}

func TestService_Min_success(t *testing.T) {
	s := newTestService()

	min := s.min(1, 2)
	if min != 1 {
		t.Fail()
	}
}

func TestService_Min_fail(t *testing.T) {
	s := newTestService()

	min := s.min(1, 2)
	if min == 2 {
		t.Fail()
	}
}

func TestService_FileExist_success(t *testing.T) {
	s := &Service{}

	err := s.ExportToFile("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}

	isExist := s.fileExist("accounts.txt")
	if !isExist {
		t.Fail()
	}

	err = os.Remove("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_GetFullPath_success(t *testing.T) {
	s := &Service{}

	err := s.ExportToFile("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}

	path, err := s.getFullPath(".", "accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}

	fullpath := filepath.Dir(path) + "/accounts.txt"
	if fullpath != path {
		t.Fail()
		return
	}

	err = os.Remove("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ExportAccountsToFile(t *testing.T) {
	s := &Service{}

	path := "accounts.txt"
	err := s.exportAccountsToFile(path, "\n")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = os.Stat(path)
	if os.IsExist(err) {
		t.Error(err)
		return
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ExportPaymentsToFile(t *testing.T) {
	s := &Service{}

	path := "payments.txt"
	err := s.exportPaymentsToFile(path)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = os.Stat(path)
	if os.IsExist(err) {
		t.Error(err)
		return
	}
	
	err = os.Remove(path)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ExportFavoritesToFile(t *testing.T) {
	s := &Service{}

	path := "favorites.txt"
	err := s.exportFavoritesToFile(path)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = os.Stat(path)
	if os.IsExist(err) {
		t.Error(err)
		return
	}
	
	err = os.Remove(path)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ImportAccountsFromFile(t *testing.T) {
	s := &Service{}

	path := "accounts.txt"
	err := s.exportAccountsToFile(path, "\n")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.importAccountsFromFile(path)
	if err != nil {
		t.Fail()
		return
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ImportPaymentsFromFile(t *testing.T) {
	s := &Service{}

	path := "payments.txt"
	err := s.exportPaymentsToFile(path)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.importPaymentsFromFile(path)
	if err != nil {
		t.Fail()
		return
	}
	
	err = os.Remove(path)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ImportFavoritesFromFile(t *testing.T) {
	s := &Service{}

	path := "favorites.txt"
	err := s.exportFavoritesToFile(path)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.importFavoritesFromFile(path)
	if err != nil {
		t.Fail()
		return
	}
	
	err = os.Remove(path)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_GetDataFromFile(t *testing.T) {
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

	err = s.ExportToFile("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}

	data, err := s.getDataFromFile("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}

	result := "1;+992937452945;100|2;+992937452946;101|3;+992937452947;102|"
	if data != result {
		t.Fail()
	}

	err = os.Remove("accounts.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ParseAccountToString(t *testing.T) {
	s := &Service{}

	account, err := s.RegisterAccount("+992937452945")
	if err != nil {
		t.Error(err)
		return
	}

	data := s.parseAccountToString(account, "|")
	expected := "1;+992937452945;0|"
	if data != expected {
		t.Fail()
		return
	}
}

func TestService_ParseStringToAccount(t *testing.T) {
	s := &Service{}

	account, err := s.RegisterAccount("+992937452945")
	if err != nil {
		t.Error(err)
		return
	}
	
	result := s.parseStringToAccounts("1;+992937452945;0|", "|")[0]
	if !reflect.DeepEqual(result, account) {
		t.Fail()
	}
}

func TestService_ParsePaymentToString(t *testing.T) {
	s := &Service{}

	payment := &types.Payment {
		ID:				"1",
		AccountID:		1,
		Amount:			10,
		Category:		"auto",
		Status:			types.PaymentStatusOk,
	}

	expected := "1;1;10;auto;OK"
	result := strings.TrimSpace(s.parsePaymentToString(payment))
	
	if result != expected {
		t.Fail()
	}
}

func TestService_ParseStringToPayment(t *testing.T) {
	s := &Service{}

	expected := &types.Payment {
		ID:				"1",
		AccountID:		1,
		Amount:			10,
		Category:		"auto",
		Status:			types.PaymentStatusOk,
	}

	data := "1;1;10;auto;OK"
	result := s.parseStringToPayments(data)[0]
	
	if !reflect.DeepEqual(result, expected) {
		t.Fail()
	}
}

func TestService_ParseFavoriteToString(t *testing.T) {
	s := &Service{}

	favorite := &types.Favorite {
		ID:				"1",
		AccountID:		1,
		Name:			"auto",
		Amount:			10,
		Category:		"auto",
	}

	expected := "1;1;auto;10;auto"
	result := strings.TrimSpace(s.parseFavoriteToString(favorite))
	
	if result != expected {
		t.Fail()
	}
}

func TestService_ParseStringToFavorite(t *testing.T) {
	s := &Service{}

	expected := &types.Favorite {
		ID:				"1",
		AccountID:		1,
		Name:			"auto",
		Amount:			10,
		Category:		"auto",
	}

	data := "1;1;auto;10;auto"
	result := s.parseStringToFavorites(data)[0]
	
	if !reflect.DeepEqual(result, expected) {
		t.Fail()
	}
}

func TestService_ContainsAccount(t *testing.T) {
	s := &Service{}

	account, err := s.RegisterAccount("+992937452945")
	if err != nil {
		t.Error(err)
		return
	}

	account2, err := s.RegisterAccount("+992937452946")
	if err != nil {
		t.Error(err)
		return
	}

	accounts := []*types.Account{account, account2}
	if !s.containsAccount(account, accounts) {
		t.Fail()
	}
}

func TestService_ContainsPayment(t *testing.T) {
	s := &Service{}

	payment := &types.Payment {
		ID:				"1",
		AccountID:		1,
		Amount:			10,
		Category:		"auto",
		Status:			types.PaymentStatusOk,
	}

	payment2 := &types.Payment {
		ID:				"2",
		AccountID:		1,
		Amount:			10,
		Category:		"auto",
		Status:			types.PaymentStatusOk,
	}

	payments := []*types.Payment{payment, payment2}
	if !s.containsPayment(payment, payments) {
		t.Fail()
	}
}

func TestService_ContainsFavorite(t *testing.T) {
	s := &Service{}

	favorite := &types.Favorite {
		ID:				"1",
		AccountID:		1,
		Name:			"auto",
		Amount:			10,
		Category:		"auto",
	}

	favorite2 := &types.Favorite {
		ID:				"2",
		AccountID:		1,
		Name:			"auto",
		Amount:			10,
		Category:		"auto",
	}

	favorites := []*types.Favorite{favorite, favorite2}
	if !s.containsFavorite(favorite, favorites) {
		t.Fail()
	}
}

func TestService_FilterPayment(t *testing.T) {
	s := &Service{}

	payment := &types.Payment {
		ID:				"1",
		AccountID:		1,
		Amount:			10,
		Category:		"book",
		Status:			types.PaymentStatusOk,
	}

	payment2 := &types.Payment {
		ID:				"2",
		AccountID:		1,
		Amount:			10,
		Category:		"auto",
		Status:			types.PaymentStatusOk,
	}

	s.payments = append(s.payments, payment)
	s.payments = append(s.payments, payment2)
	if !s.filter(*payment) {
		t.Fail()
	}
}

func BenchmarkSumOfPaymentsRegular(b *testing.B) {
	s := newTestService()
	_, _, _, err := s.addAcount(defaultTestAccount)
	if err != nil {
		b.Error(err)
		return
	}

	want := types.Money(1_000_00)
	for i := 0; i < b.N; i++ {
		result := s.sumOf(s.payments)
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
	}	
}

func BenchmarkSumOfPaymentsConcurrently(b *testing.B) {
	s := newTestService()
	_, _, _, err := s.addAcount(defaultTestAccount)
	if err != nil {
		b.Error(err)
		return
	}

	want := types.Money(1_000_00)
	for i := 0; i < b.N; i++ {
		result := s.SumPayments(2)
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
	}
}

func BenchmarkFilterPaymentsRegular(b *testing.B) {
	s := newTestService()

	accounts := []*types.Account {
		{ ID: 1, Phone: "111111", Balance: 0 },
		{ ID: 2, Phone: "111111", Balance: 0 },
		{ ID: 3, Phone: "111111", Balance: 0 },
		{ ID: 4, Phone: "111111", Balance: 0 },
	}

	s.accounts = accounts

	payments := []*types.Payment {
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 4, Amount: 1, Category: "auto" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 3, Amount: 1, Category: "auto" },
	}
	
	s.payments = payments
	want := 2
	for i := 0; i < b.N; i++ {
		filtered, err := s.FilterPayments(2, 1)
		if err != nil {
			b.Error(err)
			return
		}
		if len(filtered) != want {
			b.Fatalf("invalid result, got %v, want %v", len(filtered), want)
		}
	}
}

func BenchmarkFilterPaymentsConcurrently(b *testing.B) {
	s := newTestService()

	accounts := []*types.Account {
		{ ID: 1, Phone: "111111", Balance: 0 },
		{ ID: 2, Phone: "111111", Balance: 0 },
		{ ID: 3, Phone: "111111", Balance: 0 },
		{ ID: 4, Phone: "111111", Balance: 0 },
	}

	s.accounts = accounts

	payments := []*types.Payment {
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 4, Amount: 1, Category: "auto" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 3, Amount: 1, Category: "auto" },
	}
	
	s.payments = payments
	want := 2
	for i := 0; i < b.N; i++ {
		filtered, err := s.FilterPayments(2, 3)
		if err != nil {
			b.Error(err)
			return
		}
		if len(filtered) != want {
			b.Fatalf("invalid result, got %v, want %v", len(filtered), want)
		}
	}
}

func BenchmarkFilterPaymentsByFnConcurrently(b *testing.B) {
	s := newTestService()

	payments := []*types.Payment {
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "book" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "book" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 1, Amount: 1, Category: "auto" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 4, Amount: 1, Category: "auto" },
		{ AccountID: 2, Amount: 1, Category: "auto" },
		{ AccountID: 3, Amount: 1, Category: "book" },
	}
	
	s.payments = payments
	want := 3
	for i := 0; i < b.N; i++ {
		filtered, err := s.FilterPaymentsByFn(s.filter, 3)
		if err != nil {
			b.Error(err)
			return
		}
		if len(filtered) != want {
			b.Fatalf("invalid result, got %v, want %v", len(filtered), want)
		}
	}
}

func BenchmarkSumPaymentsWithProgress(b *testing.B) {
	s := newTestService()

	count := 1_000_000
	var payments []*types.Payment
	for i := 0; i < count; i++ {
		payment := &types.Payment {
			AccountID: 1, Amount: 1, Category: "auto",
		}

		payments = append(payments, payment)
	}

	s.payments = payments
	for i := 0; i < b.N; i++ {
		ch := s.SumPaymentsWithProgress()

		total := types.Money(0)
		for i := 0; i < count; i++ {
			progress := <- ch
			total += progress.Result
		}

		want := types.Money(1_000_000)
		if total != want {
			b.Errorf("invalid result! Expected %v, got %v", want, total)
		}
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