package usersusecases

import (
	"fmt"

	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/users"
	usersRepositories "github.com/NattpkJsw/real-world-api-go/modules/users/usersRepositories"
	"github.com/NattpkJsw/real-world-api-go/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.usersRepository.InsertUser(req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	//Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	//Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}
	// sign token
	accessToken, _ := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
	})

	// refresh token
	refreshToken, _ := auth.NewAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
	})

	//Set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}
	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}

// func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
// 	// Parse token
// 	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
// 	if err != nil {
// 		return nil, err
// 	}

// 	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
// 	if err != nil {
// 		return nil, err
// 	}

// 	Find profile
// 	profile, err := u.usersRepository.GetProfile(oauth.UserId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	newClaims := &users.UserClaims{
// 		Id:     profile.Id,
// 	}

// 	accessToken, err := auth.NewAuth(
// 		auth.Access,
// 		u.cfg.Jwt(),
// 		newClaims,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	refreshToken := auth.RepeatToken(
// 		u.cfg.Jwt(),
// 		newClaims,
// 		claims.ExpiresAt.Unix(),
// 	)

// 	passport := &users.UserPassport{
// 		User: profile,
// 		Token: &users.UserToken{
// 			Id:           oauth.Id,
// 			AccessToken:  accessToken.SignToken(),
// 			RefreshToken: refreshToken,
// 		},
// 	}
// 	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
// 		return nil, err
// 	}
// 	return passport, nil
// }
