package models

type Post struct {
	PostId string `json:"post_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
