package domain

type Profile struct {
	ProfileID string  `json:"-"`
	FullName  string  `json:"full_name"`
	AccountID string  `json:"-"`
	Account   Account `json:"-"`
}

type IProfileRepository interface {
	InsertProfile(profile Profile) error
	GetProfile(filter ProfileFilter) (*Profile, error)
	UpdateFullName(profile Profile) error
}

type IProfileUsecase interface {
	GetProfile(profile Profile) (*Profile, error)
	UpdateProfile(profile Profile) error
}

type ProfileFilter struct {
	ProfileID string
	AccountID string
}
