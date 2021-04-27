package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/rgsgit/wallet/pkg/types"
)

type testService struct {
	*Service
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{
			amount:   10_00,
			category: "auto",
		},
		{
			amount:   10_00,
			category: "food",
		},
		{
			amount:   10_00,
			category: "auto",
		},
	},
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
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

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount("")
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("Reject() can't create payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func TestService_FindAccountByID_exist(t *testing.T) {
	svc := Service{
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "992900000001",
				Balance: 300_00,
			},
			{
				ID:      2,
				Phone:   "992900000002",
				Balance: 700_00,
			},
		},
	}

	expected := &types.Account{
		ID:      1,
		Phone:   "992900000001",
		Balance: 300_00,
	}

	result, _ := svc.FindAccountByID(1)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Invalid Result: Excpected: %v, Got: %v ", expected, result)
	}
}
func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := Service{
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "992900000001",
				Balance: 300_00,
			},
			{
				ID:      2,
				Phone:   "992900000002",
				Balance: 700_00,
			},
		},
	}

	result, err := svc.FindAccountByID(3)
	if err == nil {
		t.Errorf("Invalid Got: %v, want: %v", *result, nil)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	svc := Service{
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "992900000001",
				Balance: 300_00,
			},
			{
				ID:      2,
				Phone:   "992900000002",
				Balance: 700_00,
			},
		},

		payments: []*types.Payment{
			{
				ID:        "1",
				AccountID: 2,
				Amount:    100_00,
				Status:    types.PaymentStatusInProgress,
			},
		},
	}

	/*payment := &types.Payment{
		ID:        "1",
		AccountID: 2,
		Amount:    100_00,
		Status:    types.PaymentStatusInProgress,
	}*/

	err := svc.Reject("1")

	if err != nil {
		t.Errorf("Invalid Result: Err0: %v", err)
	}
}

func TestService_FindPaymentByID_notFound(t *testing.T) {
	svc := Service{
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "992900000001",
				Balance: 300_00,
			},
			{
				ID:      2,
				Phone:   "992900000002",
				Balance: 700_00,
			},
		},

		payments: []*types.Payment{
			{
				ID:        "1",
				AccountID: 2,
				Amount:    100_00,
				Status:    types.PaymentStatusInProgress,
			},
		},
	}

	err := svc.Reject("2")

	if err == nil {
		t.Errorf("Invalid Result: Got: %v, Want:%v", nil, ErrPaymentNotFound)
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	account, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("Reject() can't register account. Error = %v", err)
	}
	payment := payments[0]
	payment, err = s.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Reject() can't create payment. Error = %v", err)
	}
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	newPayment, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	got, err := s.FindPaymentByID(newPayment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}
	if !reflect.DeepEqual(newPayment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
	}

}

func TestService_Reject_notFound(t *testing.T) {
	service := newTestService()
	account, err := service.addAccountWithBalance("+992000000001", 100_00)
	if err != nil {
		t.Errorf("Reject() can't register account, error = %v", err)
	}

	err = service.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("Reject() can't deposit account, error = %v", err)
	}
	payment, err := service.Pay(1, 10_00, "auto")
	if err != nil {
		t.Errorf("Reject() can't create payment, error = %v", err)
	}
	err = service.Reject(payment.ID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	_, err = service.FindPaymentByID(uuid.New().String())
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): wrong must return ErrPaymentNotFound, returned = %v", err)

	}
}

func TestService_Reject_success(t *testing.T) {
	service := newTestService()
	account, payments, err := service.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("Reject() can't register account. error = %v", err)
	}
	payment := payments[0]
	payment, err = service.Pay(payment.AccountID, 10_00, "auto")
	if err != nil {
		t.Errorf("Reject() can't create payment, error = %v", err)
	}
	err = service.Reject(payment.ID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	savedPayment, err := service.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	savedAccount, err := service.FindAccountByID(savedPayment.AccountID)
	if err != nil {
		t.Errorf("FindAccountByID(): error = %v", err)
		return
	}
	if savedAccount.Balance != account.Balance {
		t.Errorf("Balance didn't change: error = %v", savedAccount)
	}
}

func TestServive_GetFavoriteByID_success(t *testing.T) {
	s := newTestService()
	account, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("Reject() can't register account. Error = %v", err)
	}
	payment := payments[0]
	payment, err = s.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Reject() can't create payment. Error = %v", err)
		return
	}
	favorite, err := s.FavoritePayment(payment.ID, "First payment")
	if err != nil {
		t.Errorf("Error: %v ", err)
		return
	}
	got, err := s.GetFavoriteByID(favorite.ID)
	if err != nil {
		t.Errorf("Error: %v ", err)
		return
	}
	if !reflect.DeepEqual(favorite, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
	}
}
func TestServive_GetFavoriteByID_notFound(t *testing.T) {
	s := newTestService()
	account, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("Reject() can't register account. Error = %v", err)
	}
	payment := payments[0]
	payment, err = s.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Reject() can't create payment. Error = %v", err)
		return
	}
	_, err = s.FavoritePayment(payment.ID, "First payment")
	if err != nil {
		t.Errorf("Error: %v ", err)
		return
	}
	_, err = s.GetFavoriteByID(uuid.New().String())
	if err == nil {
		t.Errorf("Error: %v ", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	account, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("Reject() can't register account. Error = %v", err)
	}
	payment := payments[0]
	payment, err = s.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Reject() can't create payment. Error = %v", err)
		return
	}
	favorite, err := s.FavoritePayment(payment.ID, "First payment")
	if err != nil {
		t.Errorf("Error: %v ", err)
		return
	}

	newPayment, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	got, err := s.FindPaymentByID(newPayment.ID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if !reflect.DeepEqual(newPayment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
	}
}

func Transactions(s *testService) {
	s.RegisterAccount("1111")
	s.Deposit(1, 500)
	s.Pay(1, 10, "food")
	s.Pay(1, 10, "phone")
	s.Pay(1, 15, "cafe")
	s.Pay(1, 25, "auto")
	s.Pay(1, 30, "restaurant")
	s.Pay(1, 50, "auto")
	s.Pay(1, 60, "bank")
	pay, _ := s.Pay(1, 50, "bank")
	s.FavoritePayment(pay.ID, "50_for_bank")

	s.RegisterAccount("2222")
	s.Deposit(2, 200)
	pay1, _ := s.Pay(2, 40, "phone")
	s.FavoritePayment(pay1.ID, "40_for_phone")

	s.RegisterAccount("3333")
	s.Deposit(3, 300)
	s.Pay(3, 36, "auto")
	s.Pay(3, 12, "food")
	pay2, _ := s.Pay(3, 25, "phone")
	s.FavoritePayment(pay2.ID, "25_for_phone")
}

func Transactions2(s *testService) {
	s.RegisterAccount("1111")
	s.Deposit(1, 500)
	s.Pay(1, 10, "food")
	s.Pay(1, 10, "phone")
	s.Pay(1, 15, "bank")
	s.Pay(1, 25, "auto")
	s.Pay(1, 30, "restaurant")
	s.Pay(1, 50, "auto")
	s.Pay(1, 60, "bank")
	s.Pay(1, 50, "bank")

	s.RegisterAccount("2222")
	s.Deposit(2, 200)
	s.Pay(2, 40, "phone")

	s.RegisterAccount("3333")
	s.Deposit(3, 300)
	s.Pay(3, 36, "auto")
	s.Pay(3, 12, "food")
	s.Pay(3, 25, "phone")
}

func TestService_ExportToFile_success(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.ExportToFile("file.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ExportToFile_notFound(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.ExportToFile("")
	if err == nil {
		t.Error(err)
		return
	}
}

func TestService_ImportFromFile_success(t *testing.T) {
	s := newTestService()
	s.RegisterAccount("1111")
	s.Deposit(1, 500)
	pay, _ := s.Pay(1, 100, "phone")
	s.FavoritePayment(pay.ID, "my_phone")

	err := s.ImportFromFile("file.txt")
	if err != nil {
		t.Error(err)
		return
	}
}
func TestService_ImportFromFile_noSuccess(t *testing.T) {
	s := newTestService()

	err := s.ImportFromFile("")
	if err == nil {
		t.Error(err)
		return
	}
}

func TestService_Export_success(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	Transactions(s)

	err = s.Export("../../data")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Import_success(t *testing.T) {
	s := newTestService()
	s.RegisterAccount("1111")
	s.Deposit(1, 500)
	pay, _ := s.Pay(1, 100, "phone")
	s.FavoritePayment(pay.ID, "my_phone")

	err := s.Import("../../data")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Export_noSuccess(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	Transactions(s)

	err = s.Export("../data")
	if err == nil {
		t.Error(err)
		return
	}
}

func TestService_Import_noSuccess(t *testing.T) {
	s := newTestService()

	err := s.Import("../nodir")
	if err == nil {
		t.Error(err)
		return
	}
}

func BenchmarkSumPayments(b *testing.B) {
	s := newTestService()
	Transactions(s)
	want := types.Money(363)
	for i := 0; i < b.N; i++ {
		result := s.SumPayments(3)
		if result != want {
			b.Fatalf("INVALID: result_we_got %v, result_we_want %v", result, want)
		}
	}
}


func TestService_ExportAccountHistory_success(t *testing.T) {
	s := newTestService()
	Transactions(s)
	_, err := s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
}

func TestService_ExportAccountHistory_notSuccess1(t *testing.T) {
	s := newTestService()
	s.RegisterAccount("")
	s.Deposit(0, 0)
	s.Pay(0, 0, "")
	_, err := s.ExportAccountHistory(1)
	if err == nil {
		t.Error(err)
	}
}
func TestService_ExportAccountHistory_notSuccess2(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	anotherID := s.nextAccountID + 1
	_, err = s.FindAccountByID(anotherID)
	if err == nil {
		t.Error("ExportAccountHistory(): must return error, returned nil")
	}

	_, err = s.ExportAccountHistory(3)
	if err == nil {
		t.Error(err)
	}
	if err != ErrAccountNotFound {
		t.Errorf("ExportAccountHistory(): must return ErrAccountNotFound, returned = %v", err)
		return
	}

}

func TestService_HistoryToFiles_success1(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payments, err := s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
	err = s.HistoryToFiles(payments, "data", 3)
	if err != nil {
		t.Error(err)
	}
}
func TestService_HistoryToFiles_Success2(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payments, err := s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
	err = s.HistoryToFiles(payments, "data", 12)
	if err != nil {
		t.Error(err)
	}
}

func TestService_HistoryToFiles_notSuccess1(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payment := []types.Payment{}
	err := s.HistoryToFiles(payment, "data", 12)
	if err != nil {
		t.Error(err)
	}
}
func TestService_HistoryToFiles_notSuccess2(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payment := []types.Payment{}
	err := s.HistoryToFiles(payment, "", 0)
	if err == nil {
		t.Error(err)
	}
}

func TestService_SumPayments(t *testing.T) {
	s := newTestService()
	Transactions(s)
	sum := s.SumPayments(0)
	if sum != 363 {
		t.Errorf("TestService_SumPayments(): sum=%v", sum)
	}

}

func TestService_FilterPayments_success(t *testing.T) {
	s := newTestService()
	Transactions(s)

	paymnet, err := s.FilterPayments(1, 0)
	if err != nil {
		t.Error(err)
	}

	want := 8
	result := len(paymnet)
	if !reflect.DeepEqual(result, want) {
		t.Errorf("INVALID: result_we_got %v, result_we_want %v", result, want)
		return
	}
}
func TestService_FilterPayments_not_Success(t *testing.T) {
	s := newTestService()
	Transactions(s)

	_, err := s.FilterPayments(0, 0)
	if err == nil {
		t.Error(err)
	}
}

func BenchmarkFilterPayments(b *testing.B) {
	s := newTestService()
	Transactions(s)

	for i := 0; i < b.N; i++ {
		paymnet, err := s.FilterPayments(1, 3)
		if err != nil {
			b.Error(err)
		}

		want := 8
		result := len(paymnet)
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("INVALID: result_we_got %v, result_we_want %v", result, want)
			return
		}
	}
}

func TestService_FilterPaymentsByFn(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payment, err := s.FilterPaymentsByFn(FilterCategory, 0)
	if err != nil {
		t.Error(err)
	}

	want := 2
	result := len(payment)
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("INVALID: result_we_got %v, result_we_want %v", result, want)
	}
}

func BenchmarkFilterPaymentsByFn(b *testing.B) {
	s := newTestService()
	Transactions(s)

	for i := 0; i < b.N; i++ {
		payment, err := s.FilterPaymentsByFn(FilterCategory, 10)
		if err != nil {
			b.Error(err)
		}

		want := 3
		result := len(payment)
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("INVALID: result_we_got %v, result_we_want %v", result, want)
		}
	}
}

