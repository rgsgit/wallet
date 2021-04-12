package types

//Money the minimal money unit
type Money int64

//PaymentCategory пердставляет собой категорию, в которй был соверщен платеж(авто, аптеки, рестораны и т.д)
type PaymentCategory string

//PaymentStatus представляет собой статус платежа.
type PaymentStatus string

//Предопределенные статусы платеже
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

//Payment представляет информация о платеже
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

//Phone номер телефона
type Phone string

//Account представляет информацию о счёте пользователя.
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}


