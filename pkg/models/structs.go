package models

type Post struct {
	UserID int
	ID     int
	Title  string
	Body   string
}

type Comment struct {
	PostID int
	ID     int
	Name   string
	Email  string
	Body   string
}
