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
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite payment not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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

/*type Error string

func (e Error) Error() string {
	return string(e)

}*/

//FindAccountByID поис аккаунта по ID
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var accaount *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			accaount = acc
			break
		}
	}

	if accaount == nil {
		return nil, ErrAccountNotFound
	}

	return accaount, nil
}

//FindPaymentByID поиск плотежа по ID
func (s *Service) FindPaymentByID(paymetID string) (*types.Payment, error) {
	var payment *types.Payment

	for _, pmt := range s.payments {
		if pmt.ID == paymetID {
			payment = pmt
			break
		}
	}

	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil

}

//Reject метод отмены платежа
func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if payment == nil {
		return err
	}

	if payment.Status == types.PaymentStatusFail {
		return nil
	}

	acc, err := s.FindAccountByID(payment.AccountID)
	if acc == nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	acc.Balance += payment.Amount

	return nil
}

//Repeat повторяет платёж по идинтификатору
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if payment == nil {
		return nil, err
	}

	newPayment, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if newPayment == nil {
		return nil, err
	}

	return newPayment, nil

}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) GetFavoriteByID(favoriteID string) (*types.Favorite, error) {
	var favorite *types.Favorite
	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			favorite = fav
			break
		}
	}
	if favorite == nil {
		return nil, ErrFavoriteNotFound
	}
	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.GetFavoriteByID(favoriteID)
	if err != nil {
		return nil, ErrFavoriteNotFound
	}
	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}
	return payment, nil
}
