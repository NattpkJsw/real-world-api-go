package usershandlers

import (
	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/entities"
	"github.com/NattpkJsw/real-world-api-go/modules/users"
	usersUsecases "github.com/NattpkJsw/real-world-api-go/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type userHandlersErrCode string

const (
	signUpErr          userHandlersErrCode = "users-001"
	signInErr          userHandlersErrCode = "users-002"
	refreshPassportErr userHandlersErrCode = "users-003"
	signOutErr         userHandlersErrCode = "users-004"
	getUserProfileErr  userHandlersErrCode = "users-005"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	// RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpErr),
			"email pattern is invalid",
		).Res()
	}

	// Insert
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code, //500
				string(signUpErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

// func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
// 	req := new(users.UserRefreshCredential)
// 	if err := c.BodyParser(req); err != nil {
// 		return entities.NewResponse(c).Error(
// 			fiber.ErrBadRequest.Code,
// 			string(refreshPassportErr),
// 			err.Error(),
// 		).Res()
// 	}

// 	passport, err := h.usersUsecase.RefreshPassport(req)
// 	if err != nil {
// 		return entities.NewResponse(c).Error(
// 			fiber.ErrBadRequest.Code,
// 			string(refreshPassportErr),
// 			err.Error(),
// 		).Res()
// 	}
// 	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
// }

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.Oauth)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}
	if err := h.usersUsecase.DeleteOauth(req.AccessToken); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int)

	// Get profile
	result, err := h.usersUsecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}
