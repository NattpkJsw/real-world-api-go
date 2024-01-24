package articlesusecases

import (
	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/articles"
	articlesrepositories "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesRepositories"
)

type IArticlesUsecase interface {
	GetSingleArticle(slug string, userId int) (*articles.Article, error)
	GetArticlesList(req *articles.ArticleFilter, userId int) (*articles.ArticleList, error)
	GetArticlesFeed(req *articles.ArticleFeedFilter, userId int) (*articles.ArticleList, error)
	CreateArticle(req *articles.ArticleCredential) (*articles.Article, error)
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
	articleId, err := u.articlesRepository.GetArticleIdBySlug(slug)
	if err != nil {
		return nil, err
	}

	article, err := u.articlesRepository.GetSingleArticle(articleId, userId)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (u *articlesUsecase) GetArticlesList(req *articles.ArticleFilter, userId int) (*articles.ArticleList, error) {
	articleList, count, err := u.articlesRepository.GetArticlesList(req, userId)

	articlesOut := &articles.ArticleList{
		Article:       articleList,
		ArticlesCount: count,
	}

	return articlesOut, err
}

func (u *articlesUsecase) GetArticlesFeed(req *articles.ArticleFeedFilter, userId int) (*articles.ArticleList, error) {
	input := &articles.ArticleFilter{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	articleList, count, err := u.articlesRepository.GetArticlesList(input, userId)

	articlesOut := &articles.ArticleList{
		Article:       articleList,
		ArticlesCount: count,
	}

	return articlesOut, err
}

func (u *articlesUsecase) CreateArticle(req *articles.ArticleCredential) (*articles.Article, error) {
	return u.articlesRepository.CreateArticle(req)
}
