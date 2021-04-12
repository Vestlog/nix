package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	str "github.com/vestlog/nix/pkg/structs"
)

type GormDatabase struct {
	db *gorm.DB
}

func (db *GormDatabase) SavePost(post *str.Post) error {
	if err := db.db.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func (db *GormDatabase) SaveComment(comment *str.Comment) error {
	if err := db.db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func (db *GormDatabase) CreateTables() error {
	return db.db.AutoMigrate(&str.Post{}, &str.Comment{})
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
