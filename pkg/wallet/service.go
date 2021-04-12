package wallet

import "github.com/rgsgit/wallet/pkg/types"


type Service struct{
	nextAccountID	int64
	accounts 		[]types.Account
	payments 		[]types.Payment
}

func RegisterAccount(service *Service, phone types.Phone)  {
	for _, account := range service.accounts {
		if account.Phone == phone {
			return
		}
	}
	
	service.nextAccountID++
	service.accounts = append(service.accounts, types.Account{
		ID:			service.nextAccountID,
		Phone: 		phone,
		Balance: 	0,
	})
}