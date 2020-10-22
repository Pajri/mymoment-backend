package domain

type Account struct {
	AccountID string `json:"-"`
	Password  string `json:"-"`
	Email     string `json:"email"`
	Salt      []byte `json:"-"`
}

type IAccountRepository interface {
	GetAccount(filter AccountFilter) (*Account, error)
	InsertAccount(account Account) (*Account, error)
}

type AccountFilter struct {
	Email string
}
