package payload

type CreateArticleRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type ArticleInfo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	Author    *UserInfo `json:"author"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type UpdateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleQuery struct {
	AuthorName  string `query:"author"`
	QuerySearch string `query:"query"`
}

const (
	ERROR_ARTICLE_NOT_FOUND   = "article not found"
	ERROR_ARTICLE_NOT_ALLOWED = "article is not owned by author"
	ERROR_GET_ARTICLE         = "error to get articles"
)
