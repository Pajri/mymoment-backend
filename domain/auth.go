package domain

import "github.com/pajri/personal-backend/helper"

type IAuthUsecase interface {
	Login(account Account) (*helper.JWTWrapper, error)
	SignUp(account Account, profile Profile) (*Account, *Profile, error)
	VerifyEmail(token string) error
	ResetPassword(email string) error
	ChangePassword(token, password string) error
}
