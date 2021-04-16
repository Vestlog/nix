package storage

import (
	"log"
	"testing"

	"github.com/vestlog/nix/pkg/models"
)

var (
	db  *GormDatabase
	dsn = "storage.db"
)

func prepare() {
	if db == nil {
		var err error
		db, err = CreateGormDatabase(dsn)
		if err != nil {
			log.Fatal("Could not create DB")
		}
	}
	if err := db.db.AutoMigrate(&models.Post{}); err != nil {
		log.Fatal("prepare failed:", err)
	}
	if err := db.db.AutoMigrate(&models.Comment{}); err != nil {
		log.Fatal("prepare failed:", err)
	}
}

func TestSavePost(t *testing.T) {
	prepare()
	if err := db.SavePost(&models.Post{
		UserID: 0,
		ID:     0,
		Title:  "",
		Body:   "",
	}); err != nil {
		t.Error("TestSavePost failed:", err)
	}
}

func TestGetPosts(t *testing.T) {
	prepare()
	db.SavePost(&models.Post{0, 0, "title", "body"})
	posts, err := db.GetPosts()
	if err != nil || posts == nil {
		t.Error(err)
	}
}
