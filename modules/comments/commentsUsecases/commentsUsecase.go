package commentsusecases

import (
	"github.com/NattpkJsw/real-world-api-go/config"
	articlesrepositories "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesRepositories"
	"github.com/NattpkJsw/real-world-api-go/modules/comments"
	commentsrepositories "github.com/NattpkJsw/real-world-api-go/modules/comments/commentsRepositories"
)

type ICommentUsecase interface {
	FindComments(slug string, userID int) ([]*comments.Comment, error)
	InsertComment(slug string, req *comments.CommentCredential) (*comments.Comment, error)
}

type commentUsecase struct {
	cfg                config.IConfig
	commentRepository  commentsrepositories.ICommentsRepository
	articlesRepository articlesrepositories.IArticlesRepository
}

func CommentUsecase(cfg config.IConfig, commentRepository commentsrepositories.ICommentsRepository, articlesRepository articlesrepositories.IArticlesRepository) ICommentUsecase {
	return &commentUsecase{
		cfg:                cfg,
		commentRepository:  commentRepository,
		articlesRepository: articlesRepository,
	}
}

func (u *commentUsecase) FindComments(slug string, userID int) ([]*comments.Comment, error) {
	articleID, err := u.articlesRepository.GetArticleIdBySlug(slug)
	if err != nil {
		return nil, err
	}
	return u.commentRepository.FindComments(articleID, userID)
}

func (u *commentUsecase) InsertComment(slug string, req *comments.CommentCredential) (*comments.Comment, error) {
	articleID, err := u.articlesRepository.GetArticleIdBySlug(slug)
	if err != nil {
		return nil, err
	}
	req.ArticleID = articleID
	return u.commentRepository.InsertComment(req)
}
