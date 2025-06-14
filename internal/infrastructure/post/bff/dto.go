package bff

// PostsResponse представляет ответ со списком постов
type PostsResponse struct {
	Posts []Post `json:"posts"`
}

// CreatePostResponse представляет ответ на создание поста
type CreatePostResponse struct {
	Message string `json:"message"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

type Image struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

type Post struct {
	Uuid    string `json:"uuid"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Images  []*Image `json:"images"`
}

type GetPostsResponse struct {
	Posts []*Post `json:"posts"`
	Total int64 `json:"total"`
	Page int64 `json:"page"`
	Limit int64 `json:"limit"`
}