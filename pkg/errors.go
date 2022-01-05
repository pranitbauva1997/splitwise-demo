package pkg

const (
	UserUnauthenticated = iota + 1
	UserInvalid
)

type UserError struct {
	StatusCode int
}

func (err UserError) Error() string {
	switch err.StatusCode {
	case UserUnauthenticated:
		return "user is not logged in"
	case UserInvalid:
		return "user is invalid or deleted"
	default:
		return "unrecognized error occurred"
	}
}
