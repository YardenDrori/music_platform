package auth

import "github.com/YardenDrori/music-platform/internal/user"

type registerRequest struct {
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type loginRequest struct {
	Email    *string `json:"email"`
	UserName *string `json:"username"`
	Password string  `json:"password"`
}

type authResponse struct {
	User        *user.User `json:"user"`
	AccessToken string     `json:"accessToken"`
}

type accessTokenResp struct {
	AccessToken string `json:"accessToken"`
}
