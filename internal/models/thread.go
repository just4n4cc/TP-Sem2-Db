package models

type Thread struct {
	Id      int32  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Author  string `json:"author,omitempty"`
	Forum   string `json:"forum,omitempty"`
	Message string `json:"message,omitempty"`
	Votes   int32  `json:"votes,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Created string `json:"created,omitempty"`
}

type Vote struct {
	Nickname string `json:"nickname,omitempty"`
	Voice    int32  `json:"voice,omitempty"`
}
