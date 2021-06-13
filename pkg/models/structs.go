package models

type User struct {
	ID    int
	Email string
	Name  string
}

type GoogleUser struct {
	User   *User
	UserID int
	ID     string
}

type Post struct {
	UserID int
	ID     int
	Title  string
	Body   string
}

type Comment struct {
	Post   *Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostID int
	ID     int
	Name   string
	Email  string
	Body   string
}
