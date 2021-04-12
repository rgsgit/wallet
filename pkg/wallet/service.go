package wallet

import "github.com/rgsgit/wallet/pkg/types"

type Service struct {
	nextAccountID int64
	accounts      []types.Account
	payments      []types.Payment
}

//RegisterAccount регистрация аккаунта
func (service *Service) RegisterAccount(phone types.Phone) {
	for _, account := range service.accounts {
		if account.Phone == phone {
			return
		}
	}

	service.nextAccountID++
	service.accounts = append(service.accounts, types.Account{
		ID:      service.nextAccountID,
		Phone:   phone,
		Balance: 0,
	})
}

//Deposit метод пополнение счёта
func (s *Service) Deposit(accountID int64, ammount types.Money) {
	if ammount <= 0 {
		return
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = &acc
			break
		}
	}

	if account == nil {
		return
	}

	//зачисление средств
	account.Balance += ammount

}
