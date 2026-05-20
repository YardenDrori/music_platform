package user

// ==========AUTH============
type RegisterRequest struct {
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    *string
	UserName *string
	Password string
}

type AuthResponse struct {
	User         *User
	RefreshToken string
	AccessToken  string
}
