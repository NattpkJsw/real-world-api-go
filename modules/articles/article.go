package articles

type Article struct {
	Slug           *string   `json:"slug"`
	Title          *string   `json:"title"`
	Description    *string   `json:"description"`
	Body           *string   `json:"body"`
	TagList        *[]string `json:"taglist"`
	CreatedAt      *string   `json:"created_at"`
	UpdatedAt      *string   `json:"updated_at"`
	Favorited      *bool     `json:"favorited"`
	FavoritesCount *int      `json:"favcount"`
	Author         *Author   `json:"author"`
}

type Tag struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type ArticleTag struct {
	ArticleId int `db:"article_id"`
	TagId     int `db:"tag_id"`
}

type ArticleFavorite struct {
	UserId    int `db:"user_id"`
	ArticleId int `db:"article_id"`
}

type Author struct {
	Username  *string `json:"username"`
	Bio       *string `json:"bio"`
	Image     *string `json:"image"`
	Following *bool   `json:"following"`
}
