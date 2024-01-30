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
	InsertCustomer(req *users.UserRegisterReq) (*users.User, error)
	GetPassport(req *users.UserCredential) (*users.ResponsePassport, error)
	DeleteOauth(accessToken string) error
	GetUserProfile(userId int) (*users.User, error)
	UpdateUser(user *users.UserCredentialCheck) (*users.User, error)
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

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.User, error) {
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

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.ResponsePassport, error) {
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

	// set user token
	userToken := &users.UserToken{
		User_Id:      user.Id,
		AccessToken:  accessToken.SignToken(),
		RefreshToken: refreshToken.SignToken(),
	}
	if err := u.usersRepository.InsertOauth(userToken); err != nil {
		return nil, err
	}

	//Set passport
	passport := &users.UserPassport{
		// Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
		Image:    user.Image,
		Bio:      user.Bio,
		Token:    userToken.AccessToken,
	}
	passportOutput := &users.ResponsePassport{
		User: *passport,
	}
	return passportOutput, nil
}

func (u *usersUsecase) DeleteOauth(accessToken string) error {
	if err := u.usersRepository.DeleteOauth(accessToken); err != nil {
		return err
	}
	return nil
}

func (u *usersUsecase) GetUserProfile(userId int) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (u *usersUsecase) UpdateUser(user *users.UserCredentialCheck) (*users.User, error) {
	updatedUser, err := u.usersRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
