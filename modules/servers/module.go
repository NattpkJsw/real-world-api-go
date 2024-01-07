package servers

import (
	middlewareshandlers "github.com/NattpkJsw/real-world-api-go/modules/middlewares/middlewaresHandlers"
	middlewaresrepositories "github.com/NattpkJsw/real-world-api-go/modules/middlewares/middlewaresRepositories"
	middlewaresusecases "github.com/NattpkJsw/real-world-api-go/modules/middlewares/middlewaresUsecases"
	monitorhandlers "github.com/NattpkJsw/real-world-api-go/modules/monitor/monitorHandlers"
	usershandlers "github.com/NattpkJsw/real-world-api-go/modules/users/usersHandlers"
	usersrepositories "github.com/NattpkJsw/real-world-api-go/modules/users/usersRepositories"
	usersusecases "github.com/NattpkJsw/real-world-api-go/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModulefactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	middle middlewareshandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, m middlewareshandlers.IMiddlewaresHandler) IModulefactory {
	return &moduleFactory{
		router: r,
		server: s,
		middle: m,
	}
}

func InitMiddlewares(s *server) middlewareshandlers.IMiddlewaresHandler {
	repository := middlewaresrepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresusecases.MiddlewaresUsecase(repository)
	return middlewareshandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorhandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersrepositories.UsersRepository(m.server.db)
	usecase := usersusecases.UsersUsecase(m.server.cfg, repository)
	handler := usershandlers.UsersHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")
	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)

}
