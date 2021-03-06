package wallet

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

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

	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			return fav, nil
		}
	}
	return nil, ErrFavoriteNotFound
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

//ExportToFile экспортирует аккаунт в файл
func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Print(err)
		}
	}()

	for _, account := range s.accounts {
		str := strconv.FormatInt(int64(account.ID), 10) + (";") + (string(account.Phone)) + (";") + (strconv.FormatInt(int64(account.Balance), 10)) + ("|")
		_, err = file.Write([]byte(str))
		if err != nil {
			return err
		}
	}

	return nil
}

//ImportToFile импортирует даные из файла
func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Print(err)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}

		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}

	data := string(content)
	log.Println("data: ", data)

	acc := strings.Split(data, "|")
	log.Println("acc: ", acc)

	for _, operation := range acc {

		strAcc := strings.Split(operation, ";")
		log.Println("strAcc:", strAcc)
		if len(strAcc) <= 0 {
			return errors.New("Nil")
		}

		if strAcc[0] == "" {
			continue
			//return errors.New("Nil str")
		}
		id, err := strconv.ParseInt(strAcc[0], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}

		phone := types.Phone(strAcc[1])

		balance, err := strconv.ParseInt(strAcc[2], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}

		account := &types.Account{
			ID:      id,
			Phone:   phone,
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, account)
		log.Print(account)
	}
	return nil
}

//Export экспортирует все в accounts.dump, payments.dump and favorites.dump
func (s *Service) Export(dir string) error {
	if s.accounts != nil && len(s.accounts) > 0 {
		accDir, err := filepath.Abs(dir)
		if err != nil {
			log.Print(err)
			return err
		}

		accData := make([]byte, 0)

		for _, account := range s.accounts {
			str := (strconv.FormatInt(int64(account.ID), 10) + (";") +
				string(account.Phone) + (";") +
				strconv.FormatInt(int64(account.Balance), 10) + ("\n"))

			accData = append(accData, []byte(str)...)
		}
		err = os.WriteFile(accDir+"/accounts.dump", accData, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	if s.payments != nil && len(s.payments) > 0 {
		payDir, err := filepath.Abs(dir)
		if err != nil {
			log.Print(err)
			return err
		}

		payData := make([]byte, 0)

		for _, payment := range s.payments {
			str := string(payment.ID) + (";") +
				strconv.FormatInt(int64(payment.AccountID), 10) + (";") +
				strconv.FormatInt(int64(payment.Amount), 10) + (";") +
				string(payment.Category) + (";") +
				string(payment.Status) + ("\n")

			payData = append(payData, []byte(str)...)
		}
		err = os.WriteFile(payDir+"/payments.dump", payData, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	if s.favorites != nil && len(s.favorites) > 0 {
		favDir, err := filepath.Abs(dir)
		if err != nil {
			log.Print(err)
			return err
		}

		favData := make([]byte, 0)

		for _, favorite := range s.favorites {
			str := string(favorite.ID) + (";") +
				strconv.FormatInt(int64(favorite.AccountID), 10) + (";") +
				string(favorite.Name) + (";") +
				strconv.FormatInt(int64(favorite.Amount), 10) + (";") +
				string(favorite.Category) + ("\n")

			favData = append(favData, []byte(str)...)
		}
		err = os.WriteFile(favDir+"/favorites.dump", favData, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}

//Import импортирует данные из accounts.dump, payments.dump and favorites.dump

func (s *Service) Import(dir string) error {
	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Print(err)
		return err
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return err
	}

	accFile, err1 := os.ReadFile(dir + "/accounts.dump")
	if err1 == nil {

		accData := string(accFile)

		accSlice := strings.Split(accData, "\n")
		log.Print("accounts : ", accSlice)

		for _, accOperation := range accSlice {

			if len(accOperation) == 0 {
				break
			}
			accStr := strings.Split(accOperation, ";")
			log.Println("accStr:", accStr)

			id, err := strconv.ParseInt(accStr[0], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}
			phone := types.Phone(accStr[1])
			balance, err := strconv.ParseInt(accStr[2], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}

			accFind, _ := s.FindAccountByID(id)
			if accFind != nil {
				accFind.Phone = phone
				accFind.Balance = types.Money(balance)
			} else {
				s.nextAccountID++
				account := &types.Account{
					ID:      id,
					Phone:   phone,
					Balance: types.Money(balance),
				}
				s.accounts = append(s.accounts, account)
				log.Print(account)
			}
		}
	} else {
		log.Print(err1)
	}

	payFile, err2 := os.ReadFile(dir + "/payments.dump")
	if err2 == nil {

		payData := string(payFile)

		paySlice := strings.Split(payData, "\n")
		log.Print("paySlice : ", paySlice)

		for _, payOperation := range paySlice {

			if len(payOperation) == 0 {
				break
			}
			payStr := strings.Split(payOperation, ";")
			log.Println("payStr:", payStr)

			id := payStr[0]
			accountID, err := strconv.ParseInt(payStr[1], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}
			amount, err := strconv.ParseInt(payStr[2], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}
			category := types.PaymentCategory(payStr[3])
			status := types.PaymentStatus(payStr[4])

			payAcc, _ := s.FindPaymentByID(id)
			if payAcc != nil {
				payAcc.AccountID = accountID
				payAcc.Amount = types.Money(amount)
				payAcc.Category = category
				payAcc.Status = status
			} else {
				payment := &types.Payment{
					ID:        id,
					AccountID: accountID,
					Amount:    types.Money(amount),
					Category:  category,
					Status:    status,
				}
				s.payments = append(s.payments, payment)
				log.Print(payment)
			}
		}
	} else {
		log.Print(err2)
	}

	favFile, err3 := os.ReadFile(dir + "/favorites.dump")
	if err3 == nil {

		favData := string(favFile)

		favSlice := strings.Split(favData, "\n")
		log.Print("favSlice : ", favSlice)

		for _, favOperation := range favSlice {

			if len(favOperation) == 0 {
				break
			}
			favStr := strings.Split(favOperation, ";")
			log.Println("favStr:", favStr)

			id := favStr[0]
			accountID, err := strconv.ParseInt(favStr[1], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}
			name := favStr[2]
			amount, err := strconv.ParseInt(favStr[3], 10, 64)
			if err != nil {
				log.Print(err)
				return err
			}
			category := types.PaymentCategory(favStr[4])

			favAcc, _ := s.GetFavoriteByID(id)
			if favAcc != nil {
				favAcc.AccountID = accountID
				favAcc.Name = name
				favAcc.Amount = types.Money(amount)
				favAcc.Category = category
			} else {
				favorite := &types.Favorite{
					ID:        id,
					AccountID: accountID,
					Name:      name,
					Amount:    types.Money(amount),
					Category:  category,
				}
				s.favorites = append(s.favorites, favorite)
				log.Print(favorite)
			}
		}
	} else {
		log.Println(err3)
	}

	return nil

}

// SumPayments суммирует платежы
func (s *Service) SumPayments(goroutines int) types.Money {

	if goroutines < 1 {
		goroutines = 1
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	num := len(s.payments)/goroutines + 1
	sum := types.Money(0)

	for i := 0; i < goroutines; i++ {

		wg.Add(1)
		total := types.Money(0)

		func(val int) {
			defer wg.Done()
			lowIndex := val * num
			highIndex := (val * num) + num

			for j := lowIndex; j < highIndex; j++ {
				if j > len(s.payments)-1 {
					break
				}
				total += s.payments[j].Amount
			}
			mu.Lock()
			defer mu.Unlock()
			sum += total
		}(i)
	}

	wg.Wait()
	return sum
}

//ExportAccountHistory вытаскывает все платежи конктретного акаунта.
func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {

	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	payments := []types.Payment{}
	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			payments = append(payments, *payment)
		}
	}

	if len(payments) <= 0 || payments == nil {
		return nil, ErrPaymentNotFound
	}

	return payments, nil
}

//HistoryToFiles сохранение всех данных в файл.
func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {

	_, cerr := os.Stat(dir)
	if os.IsNotExist(cerr) {
		cerr = os.Mkdir(dir, 0777)
	}
	if cerr != nil {
		return cerr
	}

	if len(payments) == 0 || payments == nil {
		return nil
	}

	data := make([]byte, 0)

	if len(payments) > 0 && len(payments) <= records {
		for _, payment := range payments {
			text := []byte(
				string(payment.ID) + ";" +
					strconv.FormatInt(int64(payment.AccountID), 10) + ";" +
					strconv.FormatInt(int64(payment.Amount), 10) + ";" +
					string(payment.Category) + ";" +
					string(payment.Status) + "\n")

			data = append(data, text...)
		}

		path := dir + "/payments.dump"
		err := os.WriteFile(path, data, 0777)
		if err != nil {
			log.Print(err)
			return err
		}
	} else {
		for i, payment := range payments {

			text := []byte(
				string(payment.ID) + ";" +
					strconv.FormatInt(int64(payment.AccountID), 10) + ";" +
					strconv.FormatInt(int64(payment.Amount), 10) + ";" +
					string(payment.Category) + ";" +
					string(payment.Status) + "\n")

			data = append(data, text...)

			if (i+1)%records == 0 || i == len(payments)-1 {

				path := dir + "/payments" + strconv.Itoa((i/records)+1) + ".dump"
				err := os.WriteFile(path, data, 0777)
				if err != nil {
					log.Print(err)
					return err
				}
				data = nil
			}
		}
	}
	return nil
}

//FilterPayments отфилтровывает плотежи по accountID.
func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {

	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	if goroutines < 1 {
		goroutines = 1
	}

	num := len(s.payments)/goroutines + 1

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	payments := []types.Payment{}

	for i := 0; i < goroutines; i++ {

		wg.Add(1)
		partOfPayment := []types.Payment{}

		go func(val int) {
			defer wg.Done()
			lowIndex := val * num
			highIndex := (val * num) + num

			for j := lowIndex; j < highIndex; j++ {
				if j > len(s.payments)-1 {
					break
				}
				if s.payments[j].AccountID == accountID {
					partOfPayment = append(partOfPayment, *s.payments[j])
				}
			}
			mu.Lock()
			defer mu.Unlock()
			payments = append(payments, partOfPayment...)
		}(i)
	}

	wg.Wait()
	return payments, nil
}

//FilterPaymentsByFn - filters out payments by any function.
func (s *Service) FilterPaymentsByFn(
	filter func(payment types.Payment) bool, goroutines int) ([]types.Payment, error) {

	if goroutines < 1 {
		goroutines = 1
	}

	num := len(s.payments)/goroutines + 1

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	payments := []types.Payment{}

	for i := 0; i < goroutines; i++ {

		wg.Add(1)
		partOfPayment := []types.Payment{}

		go func(val int) {
			defer wg.Done()
			lowIndex := val * num
			highIndex := (val * num) + num

			for j := lowIndex; j < highIndex; j++ {
				if j > len(s.payments)-1 {
					break
				}
				if filter(*s.payments[j]) {
					partOfPayment = append(partOfPayment, *s.payments[j])
				}
			}
			mu.Lock()
			defer mu.Unlock()
			payments = append(payments, partOfPayment...)
		}(i)
	}

	wg.Wait()
	return payments, nil
}

func FilterCategory(payment types.Payment) bool {
	return payment.Category == "bank"
}

//SumPaymentsWithProgress разделяет в соответствие с заданным размером и суммирует все платежи в отдельных горутинах
func (s *Service) SumPaymentsWithProgress() <-chan types.Progress {

	size := 100_000

	data := []types.Money{0}
	for _, payment := range s.payments {
		data = append(data, payment.Amount)
	}

	goroutines := 1 + len(data)/size

	if goroutines <= 1 {
		goroutines = 1
	}

	channels := make([]<-chan types.Progress, goroutines)

	for i := 0; i < goroutines; i++ {

		lowIndex := i * size
		highIndex := (i + 1) * size

		if highIndex > len(data) {
			highIndex = len(data)
		}

		ch := make(chan types.Progress)
		go func(ch chan<- types.Progress, data []types.Money) {
			defer close(ch)
			sum := types.Money(0)
			for _, v := range data {
				sum += v
			}
			ch <- types.Progress{
				Part:   len(data),
				Result: sum,
			}
		}(ch, data[lowIndex:highIndex])
		channels[i] = ch
	}
	return Merge(channels)
}

//Merge возвращает канал с сообщении из всех переданных каналов
func Merge(channels []<-chan types.Progress) <-chan types.Progress {
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	merged := make(chan types.Progress)

	for _, ch := range channels {
		go func(ch <-chan types.Progress) {
			defer wg.Done()
			for val := range ch {
				merged <- val
			}
		}(ch)
	}
	go func() {
		defer close(merged)
		wg.Wait()
	}()
	return merged
}
