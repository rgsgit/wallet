package wallet

import "github.com/rgsgit/wallet/pkg/types"

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

//RegisterAccount регистрация аккаунта
func RegisterAccount(service *Service, phone types.Phone) {
	for _, account := range service.accounts {
		if account.Phone == phone {
			return
		}
	}

	service.nextAccountID++
	service.accounts = append(service.accounts, &types.Account{
		ID:      service.nextAccountID,
		Phone:   phone,
		Balance: 0,
	})
}

//RegisterAccount метод регистрация аккаунта
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, Error("phone alredy registred")
		}
	}

	s.nextAccountID++

	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}

	s.accounts = append(s.accounts, account)

	return account, nil
}

//Deposit метод пополнение счёта
func (s *Service) Deposit(accountID int64, ammount types.Money) {
	if ammount <= 0 {
		return
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return
	}

	//зачисление средств
	account.Balance += ammount

}

type Error string

func (e Error) Error() string {
	return string(e)

}
