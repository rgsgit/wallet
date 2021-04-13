package wallet

import (
	"fmt"
	"reflect"
	"testing"

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
			category: "auto",
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

	/*payment := &types.Payment{
		ID:        "1",
		AccountID: 2,
		Amount:    100_00,
		Status:    types.PaymentStatusInProgress,
	}*/

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
