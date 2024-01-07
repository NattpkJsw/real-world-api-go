package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int     `db:"id" json:"id"`
	Email    string  `db:"email" json:"email"`
	Username string  `db:"username" json:"username"`
	Image    *string `db:"image" json:"image"`
	Bio      *string `db:"bio" json:"bio"`
}

type UserRegisterReq struct {
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type UserPassport struct {
	Id       int     `db:"id" json:"id"`
	Email    string  `db:"email" json:"email"`
	Username string  `db:"username" json:"username"`
	Image    *string `db:"image" json:"image"`
	Bio      *string `db:"bio" json:"bio"`
	Token    string  `db:"access_token" json:"token"`
}

type UserToken struct {
	Id           string `db:"id" json:"id"`
	User_Id      int    `db:"user_id" json:"user_id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type UserCredential struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserCredentialCheck struct {
	Id       int     `db:"id"`
	Email    string  `db:"email"`
	Password string  `db:"password"`
	Username string  `db:"username"`
	Image    *string `db:"image"`
	Bio      *string `db:"bio"`
}

type UserClaims struct {
	Id int `db:"id" json:"id"`
}

type UserRefreshCredential struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type Oauth struct {
	AccessToken string `json:"access_token" form:"access_token"`
}

// type UserRemoveCredential struct {
// 	OauthId string `json:"oauth_id" form:"oauth_id"`
// }

func (obj *UserRegisterReq) BcryptHashing() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hash password failed: %v", err)

	}
	obj.Password = string(hashedPassword)
	return nil
}

func (obj *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
}
