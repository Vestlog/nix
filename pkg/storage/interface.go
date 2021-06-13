package storage

import "github.com/vestlog/nix/pkg/models"

type Database interface {
	SaveUser(user *models.User) error
	GetUser(id string) (*models.User, error)

	SaveGoogleUser(user *models.GoogleUser) error
	GetGoogleUser(id string) (*models.GoogleUser, error)

	GetPosts() ([]models.Post, error)
	GetPost(key string) (*models.Post, error)
	SavePost(post *models.Post) error
	UpdatePost(post *models.Post) error
	DeletePost(postid string) error

	GetComments() ([]models.Comment, error)
	GetComment(key string) (*models.Comment, error)
	SaveComment(comment *models.Comment) error
	GetCommentsPostID(postid string) ([]models.Comment, error)
	CreateTables() error
}
