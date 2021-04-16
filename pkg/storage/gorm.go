package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/vestlog/nix/pkg/models"
)

type GormDatabase struct {
	db *gorm.DB
}

func (db *GormDatabase) SavePost(post *models.Post) error {
	if err := db.db.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func (db *GormDatabase) SaveComment(comment *models.Comment) error {
	if err := db.db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func (db *GormDatabase) CreateTables() error {
	return db.db.AutoMigrate(&models.Post{}, &models.Comment{})
}

func (db *GormDatabase) GetPost(key string) (*models.Post, error) {
	dest := &models.Post{}
	if err := db.db.Where("ID = ?", key).First(dest).Error; err != nil {
		return nil, err
	}
	return dest, nil
}

func (db *GormDatabase) GetComment(key string) (*models.Comment, error) {
	dest := &models.Comment{}
	if err := db.db.Where("ID = ?", key).First(dest).Error; err != nil {
		return nil, err
	}
	return dest, nil
}

func (db *GormDatabase) GetPosts() ([]models.Post, error) {
	data := make([]models.Post, 0)
	if err := db.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (db *GormDatabase) GetComments() ([]models.Comment, error) {
	data := make([]models.Comment, 0)
	if err := db.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
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
