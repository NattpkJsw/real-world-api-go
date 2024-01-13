package articlesusecases

import (
	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/articles"
	articlesrepositories "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesRepositories"
)

type IArticlesUsecase interface {
	GetSingleArticle(slug string, userId int) (*articles.Article, error)
}

type articlesUsecase struct {
	cfg                config.IConfig
	articlesRepository articlesrepositories.IArticlesRepository
}

func ArticlesUsecase(cfg config.IConfig, articlesRepository articlesrepositories.IArticlesRepository) IArticlesUsecase {
	return &articlesUsecase{
		cfg:                cfg,
		articlesRepository: articlesRepository,
	}
}

func (u *articlesUsecase) GetSingleArticle(slug string, userId int) (*articles.Article, error) {
	article, err := u.articlesRepository.GetSingleArticle(slug, userId)
	if err != nil {
		return nil, err
	}
	return article, nil
}
