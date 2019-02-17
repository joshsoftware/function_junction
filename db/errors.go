package db

import "errors"

var (
	ErrCategoryNotExist = errors.New("Category does not exist in db")
	ErrEventNotExist = errors.New("Event does not exist in db")
)
