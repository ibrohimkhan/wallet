package types

// Money in cents
type Money int64

// PaymentCategory - category of the payment
type PaymentCategory string

// PaymentStatus - status of the payment
type PaymentStatus string

// Predefined status values
const (
	PaymentStatusOk			PaymentStatus = "OK"
	PaymentStatusFail		PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment info
type Payment struct {
	ID			string
	Amount		Money
	Category 	PaymentCategory
	Status 		PaymentStatus
}

// Phone number
type Phone string

// Account info
type Account struct {
	ID		int64
	Phone	Phone
	Balance	Money
}