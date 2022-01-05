package store

type UserResponse struct {
	Id        int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
}
