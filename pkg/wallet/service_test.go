package wallet

import (
	"reflect"
	"testing"

	"github.com/rgsgit/wallet/pkg/types"
)

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
