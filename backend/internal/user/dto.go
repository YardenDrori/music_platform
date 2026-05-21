package user

// ==========AUTH============
type registerRequest struct {
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type loginRequest struct {
	Email    *string
	UserName *string
	Password string
}

type authResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"accessToken"`
}
