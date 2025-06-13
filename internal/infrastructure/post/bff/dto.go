package bff

type Post struct {
	Uuid    string `json:"uuid"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Images  []string `json:"images"`
}

type CreatePostPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}