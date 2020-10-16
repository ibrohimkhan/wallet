package main

import (
	"sync"
	"path/filepath"
	"log"
	"os"
	"fmt"
	"github.com/ibrohimkhan/wallet/v1.1.0/pkg/wallet"
)

func main() {
	//history()
}

func min(a int, b int) int  {
	if a <= b {
		return a
	}
	return b
}

func sum1(items []int) int {
	amount := 0
	for _, item := range items {
		amount += item
	}

	return amount
}

func closure() {
	count := 10
	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			log.Println(val)
		}(i)
	}
	wg.Wait()
}

func concurrently() int {
	wg := sync.WaitGroup{}
	wg.Add(2)

	mu := sync.Mutex{}
	sum := 0

	go func() {
		defer wg.Done()
		val := 0
		for i  :=  0; i < 1_000_000; i++ {
			val++
		}
		mu.Lock()
		defer mu.Unlock()
		sum += val
	}()
		
	go func() {
		defer wg.Done()
		val := 0
		for i  :=  0; i < 1_000_000; i++ {
			val++
		}
		mu.Lock()
		defer mu.Unlock()
		sum += val
	}()

	wg.Wait()
	return sum
}

func regular() int {
	sum := 0
	for i  :=  0; i < 2_000_000; i++ {
		sum++
	}

	return sum
}

func history() {
	s := &wallet.Service{}

	createData(s)
	payments := s.GetPayments()

	path, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
		return
	}

	fullpath := path + "/data"

	err = s.HistoryToFiles(payments, fullpath, 4)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(len(payments))
}

func importData() {
	s := &wallet.Service{}

	path, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
		return
	}

	fullpath := path + "/data"
	log.Println(fullpath)

	err = s.Import(fullpath)
	if err != nil {
		log.Println(err)
		return
	}
}

func exportData() {
	s := &wallet.Service{}

	createData(s)

	path, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
		return
	}

	fullpath := path + "/data"
	log.Println(fullpath)

	err = s.Export(fullpath)
	if err != nil {
		log.Println(err)
		return
	}
}

func createData(s *wallet.Service) {
	account, err := s.RegisterAccount("+992937452945")
	if err != nil {
		log.Println(err)
		return
	}
	s.Deposit(account.ID, 100)
	s.Pay(account.ID, 15, "auto")
	payment, err := s.Pay(account.ID, 25, "auto")
	s.FavoritePayment(payment.ID, "auto")

	account, err = s.RegisterAccount("+992937452946")
	if err != nil {
		log.Println(err)
		return
	}
	s.Deposit(account.ID, 101)
	s.Pay(account.ID, 13, "auto")
	payment, err = s.Pay(account.ID, 77, "book")
	s.FavoritePayment(payment.ID, "auto")

	account, err = s.RegisterAccount("+992937452947")
	if err != nil {
		log.Println(err)
		return
	}
	s.Deposit(account.ID, 102)
	s.Pay(account.ID, 35, "book")
	s.Pay(account.ID, 33, "book")
}

func path() {
	abs, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(abs)

	path, err := filepath.Abs(abs)
	path += "/accounts.dump"
	log.Println(path)

	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(wd)

	err = os.Chdir("..")
	if err != nil {
		log.Println(err)
		return
	}

	wd, err = os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(wd)		
}

func registerAccount() {
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