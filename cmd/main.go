package main

import (
	"github.com/rgsgit/wallet/pkg/wallet"
)

func main(){

	service := &wallet.Service{}

	/*service.RegisterAccount("9920000001")
	service.Deposit(1, 500)
	pay, _ := service.Pay(1, 100, "phone")
	service.FavoritePayment(pay.ID, "my_phone")
	
	service.RegisterAccount("9920000002")
	service.Deposit(2, 1000)
	pay1, _ := service.Pay(2, 200, "auto")
	service.FavoritePayment(pay1.ID, "my_auto")
	
	service.RegisterAccount("9920000003")
	service.Deposit(3, 12000)
	pay2, _ := service.Pay(3, 300, "shop")
	service.FavoritePayment(pay2.ID, "my_shop")
	*/

	//service.Export("../data")

	service.Import("../data")
}