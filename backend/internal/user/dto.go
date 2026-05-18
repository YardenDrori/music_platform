package user

// ==========AUTH============
type RegisterRequest struct {
	Email     string
	UserName  string
	FirstName string
	LastName  string
	Password  string
}

type LoginRequest struct {
	Email    *string
	UserName *string
	Password string
}

type LoginResponse struct {
	User         *User
	RefreshToken string
	AccessToken  string
}
