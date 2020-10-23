package domain

type IAuthUsecase interface {
	Login(account Account) (string, error)
	SignUp(account Account, profile Profile) (*Account, *Profile, error)
	VerifyEmail(token string) error
	ResetPassword(email string) error
}
