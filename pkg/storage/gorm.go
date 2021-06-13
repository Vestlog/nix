package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/vestlog/nix/pkg/models"
)

type GormDatabase struct {
	DB *gorm.DB
}

func (db *GormDatabase) SaveUser(user *models.User) error {
	if err := db.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (db *GormDatabase) GetUser(id string) (*models.User, error) {
	dest := &models.User{}
	if err := db.DB.First(dest, id).Error; err != nil {
		return nil, err
	}
	return dest, nil
}

func (db *GormDatabase) SaveGoogleUser(user *models.GoogleUser) error {
	if err := db.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (db *GormDatabase) GetGoogleUser(id string) (*models.GoogleUser, error) {
	dest := &models.GoogleUser{}
	if err := db.DB.Preload("User").Where("ID = ?", id).First(dest).
		Error; err != nil {
		return nil, err
	}
	return dest, nil
}

func (db *GormDatabase) SavePost(post *models.Post) error {
	if err := db.DB.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func (db *GormDatabase) SaveComment(comment *models.Comment) error {
	if err := db.DB.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func (db *GormDatabase) GetPost(key string) (*models.Post, error) {
	dest := &models.Post{}
	if err := db.DB.Where("ID = ?", key).First(dest).Error; err != nil {
		return nil, err
	}
	return dest, nil
}

func (db *GormDatabase) GetComment(key string) (*models.Comment, error) {
	dest := &models.Comment{}
	if err := db.DB.Where("ID = ?", key).First(dest).Error; err != nil {
		return nil, err
	}
	return dest, nil
}

func (db *GormDatabase) GetPosts() ([]models.Post, error) {
	data := make([]models.Post, 0)
	if err := db.DB.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (db *GormDatabase) GetComments() ([]models.Comment, error) {
	data := make([]models.Comment, 0)
	if err := db.DB.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (db *GormDatabase) GetCommentsPostID(postid string) ([]models.Comment, error) {
	data := make([]models.Comment, 0)
	if err := db.DB.Where("post_id = ?", postid).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (db *GormDatabase) UpdatePost(post *models.Post) error {
	return db.DB.Save(post).Error
}

func (db *GormDatabase) DeletePost(postid string) error {
	return db.DB.Delete(&models.Post{}, postid).Error
}

func (db *GormDatabase) CreateTables() error {
	return db.DB.AutoMigrate(
		&models.Post{},
		&models.Comment{},
		&models.User{},
		&models.GoogleUser{},
	)
}

func CreateGormDatabase(dsn string) (*GormDatabase, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return &GormDatabase{db}, nil
}
