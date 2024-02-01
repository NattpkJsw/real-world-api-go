package comments

type Comment struct {
	Id        int                    `json:"id" db:"id"`
	CreatedAt string                 `json:"createdAt" db:"createdat"`
	UpdatedAt string                 `json:"updatedAt" db:"updatedat"`
	Body      string                 `json:"body" db:"body"`
	Author    map[string]interface{} `json:"author"`
}

type Author struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type CommentCredential struct {
	Body      string `json:"body"`
	AuthorID  int    `json:"author_id"`
	ArticleID int    `json:"article_id"`
}
