package auth

import "errors"

var (
	ErrInvalidToken  = errors.New("Отсутствует токен авторизации")
	ErrVerifyToken   = errors.New("не удалось проверить токен")
	ErrUserNotAuthor = errors.New("Пользователь не является автором поста")
)
