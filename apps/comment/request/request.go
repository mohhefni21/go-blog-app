package request

type AddCommentPayload struct {
	PostId   int    `json:"post_id"`
	ParentId int64  `json:"parent_id"`
	Content  string `json:"content"`
}

type UpdateCommentPayload struct {
	Content string `json:"content"`
}
