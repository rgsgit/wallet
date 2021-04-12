package wallet

import (
	"errors"

	"github.com/google/uuid"

	"github.com/rgsgit/wallet/pkg/types"
)

var ErrPhoneRegistred = errors.New("phone alredy registred")
var ErrAmmountMustBePositive = errors.New("ammount mus be greater then zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")

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
			return nil, ErrPhoneRegistred
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
func (s *Service) Deposit(accountID int64, ammount types.Money) error {
	if ammount <= 0 {
		return ErrAmmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	//зачисление средств
	account.Balance += ammount

	return nil

}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
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
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}

type Error string

func (e Error) Error() string {
	return string(e)

}
