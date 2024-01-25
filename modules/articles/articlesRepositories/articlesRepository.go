package articlesrepositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NattpkJsw/real-world-api-go/modules/articles"
	articlespatterns "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesPatterns"
	"github.com/jmoiron/sqlx"
)

type IArticlesRepository interface {
	GetSingleArticle(articleId int, userId int) (*articles.Article, error)
	GetArticlesList(req *articles.ArticleFilter, userId int) ([]*articles.Article, int, error)
	GetArticleIdBySlug(slug string) (int, error)
	CreateArticle(req *articles.ArticleCredential) (*articles.Article, error)
	UpdateArticle(req *articles.ArticleCredential, userID int) (*articles.Article, error)
	DeleteArticle(articleID, userID int) error
}

type articlesRepository struct {
	db *sqlx.DB
}

func ArticlesRepository(db *sqlx.DB) IArticlesRepository {
	return &articlesRepository{
		db: db,
	}
}

func (r *articlesRepository) GetSingleArticle(articleId int, userId int) (*articles.Article, error) {
	query := `
	SELECT
		to_jsonb("ar")
	FROM(
			SELECT
			"a"."slug",
			"a"."title",
			"a"."description",
			"a"."body",
			(
				SELECT coalesce(array_to_json(array_agg("t"."name")),'[]'::json)
				FROM "article_tags" "at"
				JOIN "tags" AS "t" ON "at"."tag_id" = "t"."id"
				WHERE "a"."id" = "at"."article_id"
			) AS "taglist",
			"a"."created_at",
			"a"."updated_at",
			(
				SELECT
				CASE WHEN EXISTS(
					SELECT 1
					FROM "articles" "a"
					JOIN "article_favorites" AS "af" ON "af"."article_id" = "a"."id"
					WHERE "af"."user_id" = $2
				) THEN TRUE ELSE FALSE END
			) AS "favorited",
			(
				SELECT COUNT(*)
				FROM "article_favorites" "af"
				WHERE "af"."article_id" = "a"."id"
			) AS "favoritesCount",
			(
				SELECT 
					json_build_object(
						'username', "u"."username",
						'bio', "u"."bio",
						'image', "u"."image",
						'following',
						CASE 
							WHEN EXISTS (
								SELECT 1
								FROM "user_follows" "uf"
								WHERE "a"."author_id" = "uf"."following_id"  AND "uf"."follower_id" = $2
							) THEN TRUE 
							ELSE FALSE 
						END
					)
				FROM "users" "u"
				WHERE "a"."author_id" = "u"."id"
			) AS "author"
			FROM "articles" "a"
			WHERE "a"."id" = $1
			LIMIT 1
	) AS "ar";`

	articleBytes := make([]byte, 0)
	article := new(articles.Article)

	if err := r.db.Get(&articleBytes, query, articleId, userId); err != nil {
		return nil, fmt.Errorf("get article failed: %v", err)
	}
	if err := json.Unmarshal(articleBytes, &article); err != nil {
		return nil, fmt.Errorf("unmarshal article failed: %v", err)
	}
	return article, nil
}

func (r *articlesRepository) GetArticleIdBySlug(slug string) (int, error) {
	query := `
	SELECT "a"."id"
	FROM "articles" "a"
	WHERE "a"."slug" = $1`

	var id int
	if err := r.db.Get(&id, query, slug); err != nil {
		return 0, fmt.Errorf("get articleID failed: %v", err)
	}
	return id, nil
}

func (r *articlesRepository) GetArticlesList(req *articles.ArticleFilter, userId int) ([]*articles.Article, int, error) {
	builder := articlespatterns.FindArticleBuilder(r.db, req)
	engineer := articlespatterns.FindProductEngineer(builder)

	result, err := engineer.FindArticle(userId).Result()
	count := len(result)

	return result, count, err
}

func (r *articlesRepository) CreateArticle(req *articles.ArticleCredential) (*articles.Article, error) {
	builder := articlespatterns.AddArticleBuilder(r.db, req)
	articleId, err := articlespatterns.AddArticleEngineer(builder).AddArticle()
	if err != nil {
		return nil, err
	}

	article, err := r.GetSingleArticle(articleId, req.Author)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (r *articlesRepository) UpdateArticle(req *articles.ArticleCredential, userID int) (*articles.Article, error) {
	query := `
	UPDATE "articles" SET`
	params := make(map[string]any)
	params["id"] = req.Id

	if req.Title != "" {
		query += " title = :title,"
		query += " slug = :slug,"
		params["title"] = req.Title
		params["slug"] = req.Title
	}

	if req.Body != "" {
		query += " body = :body,"
		params["body"] = req.Body
	}

	if req.Description != "" {
		query += " description = :description,"
		params["description"] = req.Description
	}

	query = query[:len(query)-1]
	query += " WHERE id = :id;"
	fmt.Println("query === ", query)
	if _, err := r.db.NamedExec(query, params); err != nil {
		return nil, fmt.Errorf("update article failed:%v", err)
	}

	return r.GetSingleArticle(req.Id, userID)
}

func (r *articlesRepository) DeleteArticle(articleID, userID int) error {
	query := `
	DELETE
	FROM "articles"
	WHERE "id" = $1 AND "author_id" = $2;`

	if _, err := r.db.ExecContext(context.Background(), query, articleID, userID); err != nil {
		return fmt.Errorf("delete article failed: %v", err)
	}
	return nil
}
