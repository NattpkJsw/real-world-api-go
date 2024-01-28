package servers

import (
	articleshandlers "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesHandlers"
	articlesrepositories "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesRepositories"
	articlesusecases "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesUsecases"
	commentshandlers "github.com/NattpkJsw/real-world-api-go/modules/comments/commentsHandlers"
	commentsrepositories "github.com/NattpkJsw/real-world-api-go/modules/comments/commentsRepositories"
	commentsusecases "github.com/NattpkJsw/real-world-api-go/modules/comments/commentsUsecases"
	"github.com/NattpkJsw/real-world-api-go/modules/middlewares"
	middlewareshandlers "github.com/NattpkJsw/real-world-api-go/modules/middlewares/middlewaresHandlers"
	middlewaresrepositories "github.com/NattpkJsw/real-world-api-go/modules/middlewares/middlewaresRepositories"
	middlewaresusecases "github.com/NattpkJsw/real-world-api-go/modules/middlewares/middlewaresUsecases"
	monitorhandlers "github.com/NattpkJsw/real-world-api-go/modules/monitor/monitorHandlers"
	profileshandlers "github.com/NattpkJsw/real-world-api-go/modules/profiles/profilesHandlers"
	profilesrepositories "github.com/NattpkJsw/real-world-api-go/modules/profiles/profilesRepositories"
	profilesusecases "github.com/NattpkJsw/real-world-api-go/modules/profiles/profilesUsecases"
	usershandlers "github.com/NattpkJsw/real-world-api-go/modules/users/usersHandlers"
	usersrepositories "github.com/NattpkJsw/real-world-api-go/modules/users/usersRepositories"
	usersusecases "github.com/NattpkJsw/real-world-api-go/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModulefactory interface {
	MonitorModule()
	UsersModule()
	ProfileModule()
	ArticleModule()
	CommentModule()
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
	router.Get("/", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.GetUserProfile)
	router.Post("/signout", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.SignOut)
	router.Put("/", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.UpdateUser)
	// router.Post("/refresh", handler.RefreshPassport)

}

func (m *moduleFactory) ProfileModule() {
	repository := profilesrepositories.ProfilesRepository(m.server.db)
	usecase := profilesusecases.ProfilesUsecase(m.server.cfg, repository)
	handler := profileshandlers.ProfileHandler(m.server.cfg, usecase)

	router := m.router.Group("/profiles")
	router.Get("/:username", m.middle.JwtAuth(string(middlewares.ReadLevel)), handler.GetProfile)
	router.Post("/:username/follow", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.FollowUser)
	router.Delete("/:username/follow", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.UnfollowUser)
}

func (m *moduleFactory) ArticleModule() {
	repository := articlesrepositories.ArticlesRepository(m.server.db)
	usecase := articlesusecases.ArticlesUsecase(m.server.cfg, repository)
	handler := articleshandlers.ArticlesHandler(m.server.cfg, usecase)

	router := m.router.Group("/articles")
	router.Get("/:slug", m.middle.JwtAuth(string(middlewares.ReadLevel)), handler.GetSingleArticle)
	router.Get("/", m.middle.JwtAuth(string(middlewares.ReadLevel)), handler.GetArticlesList)
	router.Get("/feed", m.middle.JwtAuth(string(middlewares.ReadLevel)), handler.GetArticlesFeed)
	router.Post("/", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.CreateArticle)
	router.Put("/:slug", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.UpdateArticle)
	router.Delete("/:slug", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.DeleteArticle)

}

func (m *moduleFactory) CommentModule() {
	commentRepository := commentsrepositories.CommentRepository(m.server.db)
	articleRepository := articlesrepositories.ArticlesRepository(m.server.db)
	commentUsecase := commentsusecases.CommentUsecase(m.server.cfg, commentRepository, articleRepository)
	handler := commentshandlers.CommentsHandler(m.server.cfg, commentUsecase)

	router := m.router.Group("/articles/:slug")
	router.Get("/comments", m.middle.JwtAuth(string(middlewares.ReadLevel)), handler.FindComments)
	router.Post("/comments", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.InsertComment)
	router.Delete("/comments/:id", m.middle.JwtAuth(string(middlewares.WriteLevel)), handler.DeleteComment)

}
