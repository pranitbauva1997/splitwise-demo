package store

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

type Transaction struct {
	OwedTo int64 `json:"owed_to"`
	Owes int64 `json:"owes"`
	Amount int `json:"amount"`
}
