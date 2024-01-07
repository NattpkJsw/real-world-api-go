package profilesusecases

import (
	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/profiles"
	profilesrepositories "github.com/NattpkJsw/real-world-api-go/modules/profiles/profilesRepositories"
)

type IProfilesUsecase interface {
	GetProfile(username string, curUserId int) (*profiles.Profile, error)
}

type profilesUsecase struct {
	cfg                config.IConfig
	profilesRepository profilesrepositories.IProfilesRepository
}

func ProfilesUsecase(cfg config.IConfig, profilesRepository profilesrepositories.IProfilesRepository) IProfilesUsecase {
	return &profilesUsecase{
		cfg:                cfg,
		profilesRepository: profilesRepository,
	}
}

func (u *profilesUsecase) GetProfile(username string, curUserId int) (*profiles.Profile, error) {
	profiles, err := u.profilesRepository.FindOneUserProfileByUsername(username, curUserId)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}
