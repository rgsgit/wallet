package types

//import "github.com/rgsgit/bank/v2/pkg/types"

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

//Favirite шаблон для создания платежа
type Favorite struct {
	ID        string
	AccountID int64
	Name      string
	Amount    Money
	Category  PaymentCategory
}

type Progress struct {
	Part   int
	Result Money
}
