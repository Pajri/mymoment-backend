package usecase

import "github.com/pajri/personal-backend/domain"

type ProfileUsecase struct {
	accountRepo domain.IAccountRepository
	profileRepo domain.IProfileRepository
}

func NewProfileUsecase(accountRepository domain.IAccountRepository,
	profileRepository domain.IProfileRepository) domain.IProfileUsecase {
	return &ProfileUsecase{
		accountRepo: accountRepository,
		profileRepo: profileRepository,
	}
}

func (uc ProfileUsecase) GetProfile(profile domain.Profile) (*domain.Profile, error) {
	//get account
	var accountFilter domain.AccountFilter
	accountFilter.AccountID = profile.AccountID
	storedAccount, err := uc.accountRepo.GetAccount(accountFilter)
	if err != nil {
		return nil, err
	}

	//get profile
	var profileFilter domain.ProfileFilter
	profileFilter.AccountID = storedAccount.AccountID
	storedProfile, err := uc.profileRepo.GetProfile(profileFilter)
	if err != nil {
		return nil, err
	}

	//populate profile
	profile.ProfileID = storedProfile.ProfileID
	profile.FullName = storedProfile.FullName
	profile.AccountID = storedAccount.AccountID
	profile.Account = *storedAccount

	return &profile, nil
}

func (uc ProfileUsecase) UpdateProfile(profile domain.Profile) error {
	err := uc.profileRepo.UpdateFullName(profile)
	if err != nil {
		return err
	}
	return nil
}
