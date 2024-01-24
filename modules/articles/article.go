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
	FavoritesCount *int      `json:"favoritesCount"`
	Author         *Author   `json:"author"`
}

type Author struct {
	Username  *string `json:"username"`
	Bio       *string `json:"bio"`
	Image     *string `json:"image"`
	Following *bool   `json:"following"`
}

type ArticleList struct {
	Article       []*Article `json:"article"`
	ArticlesCount int        `json:"articlesCount"`
}

type ArticleFilter struct {
	Tag       string `query:"tag"`
	Author    string `query:"author"`
	Favorited string `query:"favorited"`
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
}

type ArticleFeedFilter struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type ArticleCredential struct {
	Id          int       `json:"id"`
	Author      int       `json:"author_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	TagList     []*string `json:"tagList"`
}
