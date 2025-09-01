package examples

import (
	"fmt"
	"math/rand"
	"sync"
)

type Account struct {
	Balance int
	Mu      sync.RWMutex
}

func AccBal() {
	var account Account
	var wg sync.WaitGroup
	for range 100 {
		amount := rand.Intn(100)
		wg.Add(2)
		
		go func (amount int) {
			defer wg.Done()
			account.deposit(amount)
		} (amount)

		go func (amount int) {
			defer wg.Done()
			account.withdraw(amount)
		} (amount)
	}
	wg.Wait()
	fmt.Println(account.Balance)
}

func(a *Account) deposit(amount int) {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.Balance += amount
}

func(a *Account) withdraw(amount int) {
	a.Mu.RLock()
	defer a.Mu.RUnlock()
	a.Balance -= amount
}
