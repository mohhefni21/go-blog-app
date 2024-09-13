package request

type AddPostRequestPayload struct {
	UserId      int    `json:"user_id"`
	Title       string `json:"title"`
	Excerpt     string `json:"excerpt"`
	Content     string `json:"content"`
	Status      string `json:"status"`
	PublishedAt string `json:"published_at"`
}

type GetPostsRequestPayload struct {
	Cursor int    `query:"cursor" json:"cursor"`
	Limit  int    `query:"limit" json:"limit"`
	Search string `query:"search" json:"search"`
}

func (g *GetPostsRequestPayload) DefaultValuePagination() GetPostsRequestPayload {
	if g.Cursor < 0 {
		g.Cursor = 0
	}

	if g.Limit <= 0 {
		g.Limit = 10
	}

	return *g
}

type UpdatePostRequestPayload struct {
	Title       string `json:"title"`
	Excerpt     string `json:"excerpt"`
	Content     string `json:"content"`
	Status      string `json:"status"`
	PublishedAt string `json:"published_at"`
}
