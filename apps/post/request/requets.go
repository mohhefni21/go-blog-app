package request

type AddPostRequestPayload struct {
	UserId      int    `json:"user_id"`
	Cover       string `json:"cover"`
	Title       string `json:"title"`
	Excerpt     string `json:"excerpt"`
	Content     string `json:"content"`
	Status      string `json:"status"`
	PublishedAt string `json:"published_at"`
}
