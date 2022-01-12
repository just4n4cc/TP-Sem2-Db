package models

import "errors"

type Error struct {
	Message string `json:"message,omitempty"`
}

type Status struct {
	User   int32 `json:"user"`
	Forum  int32 `json:"forum"`
	Thread int32 `json:"thread"`
	Post   int64 `json:"post"`
}

type SortOptions struct {
	Limit int32
	//Since int64
	Since string
	Sort  string
	Desc  bool
}

var (
	AlreadyExistsError = errors.New("user already exists")
	NotFoundError      = errors.New("user not found")

	ModelFieldError = errors.New("some error with field of model")

	UnexpectedDbBehavior = errors.New("unexpected db behavior detected")
)
