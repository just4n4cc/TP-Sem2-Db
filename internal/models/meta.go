package models

type Error struct {
	Message string `json:"message,omitempty"`
}

type Status struct {
	User   int32 `json:"user,omitempty"`
	Forum  int32 `json:"forum,omitempty"`
	Thread int32 `json:"thread,omitempty"`
	Post   int64 `json:"post,omitempty"`
}
