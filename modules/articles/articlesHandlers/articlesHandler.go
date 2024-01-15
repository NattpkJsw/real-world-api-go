package articleshandlers

import (
	"net/url"
	"strings"

	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/articles"
	articlesusecases "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesUsecases"
	"github.com/NattpkJsw/real-world-api-go/modules/entities"
	"github.com/gofiber/fiber/v2"
)

type articlesHandlersErrCode string

const (
	getSingleArticleErr articlesHandlersErrCode = "article-001"
	getArticlesErr      articlesHandlersErrCode = "article-002"
)

type IArticleshandler interface {
	GetSingleArticle(c *fiber.Ctx) error
	GetArticlesList(c *fiber.Ctx) error
}

type articlesHandler struct {
	cfg             config.IConfig
	articlesUsecase articlesusecases.IArticlesUsecase
}

func ArticlesHandler(cfg config.IConfig, articlesUsecase articlesusecases.IArticlesUsecase) IArticleshandler {
	return &articlesHandler{
		cfg:             cfg,
		articlesUsecase: articlesUsecase,
	}
}

func (h *articlesHandler) GetSingleArticle(c *fiber.Ctx) error {
	pathVariable := strings.Trim(c.Params("slug"), " ")
	slug, err := url.QueryUnescape(pathVariable)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getSingleArticleErr),
			err.Error(),
		).Res()
	}
	userId := c.Locals("userId").(int)

	article, err := h.articlesUsecase.GetSingleArticle(slug, userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(getSingleArticleErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, article).Res()
}

func (h *articlesHandler) GetArticlesList(c *fiber.Ctx) error {
	req := &articles.ArticleFilter{}
	userId := c.Locals("userId").(int)

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getArticlesErr),
			err.Error(),
		).Res()
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset <= 0 {
		req.Offset = 0
	}

	articlesOut, err := h.articlesUsecase.GetArticlesList(req, userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(getArticlesErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, articlesOut).Res()

}
