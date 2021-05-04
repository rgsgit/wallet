package main

import (
	//"github.com/rgsgit/wallet/pkg/wallet"
	"log"
	"sync"
)

func main(){

	//service := &wallet.Service{}

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

	//service.Import("../data")

	data := make([]int, 1_000_000)
	for i := range data {
		data[i] = i
	}

	parts := 10
	size := len(data)/parts
	channels := make([]<-chan int, parts)
	for i :=0; i<parts;i++{
		ch := make(chan int)
		channels[i] = ch
		go func (ch chan<- int, data []int)  {
			defer close(ch)
			sum:=0
			for _, v:=range data{
				sum +=v
			}
			ch<-sum
			
		}(ch,data[i*size:(i+1)*size])
	}

	total := 0
	for value := range merge(channels){
		total+=value
	}
	log.Print(total)
}

func merge(channels []<-chan int) <-chan int{
	wg := sync.WaitGroup{}
	wg.Add(len(channels))
	merged := make(chan int)

	for _, ch := range channels{
		go func(ch<-chan int){
			defer wg.Done()
			for val:= range ch{
				merged<-val
			}
		}(ch)
	}

	go func() {
		defer close(merged)
		wg.Wait()
	}()

	return merged
}
