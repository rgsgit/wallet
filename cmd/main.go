package mamin

import "github.com/rgsgit/wallet/pkg/wallet"

func main(){
	svc := &wallet.Service{}
	wallet.RegisterAccount(svc, "+992000000001")
}