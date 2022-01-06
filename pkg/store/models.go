package store

import "errors"

type UserResponse struct {
	Id        int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
}

type AddBillRequest struct {
	CreatedBy int64 `json:"created_by"`
	Amount int `json:"amount"`
	Transactions []Transaction `json:"transactions"`
}

func (b AddBillRequest) Validate() (bool, error) {
	if b.CreatedBy == 0 {
		return false, errors.New("user_id cannot be 0")
	}
	if b.Amount == 0 {
		return false, errors.New("bill amount cannot be 0")
	}
	if b.Transactions == nil || len(b.Transactions) == 0 {
		return false, errors.New("bill cannot have 0 transactions")
	}
	sum := 0
	for _, v := range b.Transactions {
		if v.Owes == 0 || v.OwedTo == 0 {
			return false, errors.New("user_id cannot be 0")
		}
		sum = sum + v.Amount
	}

	return sum == b.Amount, nil
}

type Transaction struct {
	OwedTo int64 `json:"owed_to"`
	Owes int64 `json:"owes"`
	Amount int `json:"amount"`
}

type UserTransaction struct {
	OwedTo int64
	Owes int64
	Amount int
}
