package post

import "errors"

var (
	ErrPostNotFound = errors.New("пост не найден")
	ErrInvalidUUID  = errors.New("невалидный UUID")
)
