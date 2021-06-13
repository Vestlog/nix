package storage

import (
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/vestlog/nix/pkg/models"
)

var (
	db       *GormDatabase
	filename = "storage.db"
	options  = "?_foreign_keys=ON"
	dsn      = filename + options
)

func prepare() {
	if db == nil {
		os.Remove(filename)
		var err error
		db, err = CreateGormDatabase(dsn)
		if err != nil {
			log.Fatal("Could not create DB")
		}
	}
	if err := db.CreateTables(); err != nil {
		log.Fatal(err)
	}
}

func TestSavePost(t *testing.T) {
	prepare()
	post := &models.Post{
		UserID: 0,
		ID:     0,
		Title:  "test save post",
		Body:   "test save post",
	}
	if err := db.SavePost(post); err != nil {
		t.Error(err)
	}
	db.DB.Delete(post)
}

func TestGetPost(t *testing.T) {
	prepare()
	post := &models.Post{
		UserID: 10,
		ID:     17,
		Title:  "test get post",
		Body:   "test get post",
	}
	db.SavePost(post)
	result, err := db.GetPost("17")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(post, result) {
		t.Errorf("expected %v, got %v", post, result)
	}
	db.DB.Delete(post)
}

func TestGetPostErrNotExist(t *testing.T) {
	prepare()
	result, err := db.GetPost("17")
	if err == nil {
		t.Error("error expected, got nil and", result)
	}
}

func TestGetAllPosts(t *testing.T) {
	prepare()
	db.SavePost(&models.Post{
		UserID: 0,
		ID:     0,
		Title:  "test get all post",
		Body:   "test get all posts",
	})
	posts, err := db.GetPosts()
	if err != nil || posts == nil {
		t.Error(err)
	}
}

func TestGetCommentsPostID(t *testing.T) {
	prepare()
	db.SaveComment(&models.Comment{
		PostID: 61,
		ID:     rand.Int() % 256,
		Name:   "John Smith",
		Email:  "mail@example.com",
		Body:   "This is a comment!",
	})
	comments, err := db.GetCommentsPostID("61")
	if err != nil || comments == nil {
		t.Error(err)
	}
}

func TestUpdatePost(t *testing.T) {
	prepare()
	posts := []models.Post{
		{
			UserID: 0,
			ID:     72,
			Title:  "old title for 72",
			Body:   "old body for 72",
		},
		{
			UserID: 0,
			ID:     117,
			Title:  "old title for 117",
			Body:   "old body for 117",
		},
	}
	for _, post := range posts {
		db.SavePost(&post)
	}
	post := &models.Post{
		UserID: 200,
		ID:     72,
		Title:  "TESTNEWTITLE",
		Body:   "TESTNEWTEXT",
	}
	if err := db.UpdatePost(post); err != nil {
		t.Error(err)
	}
}

func TestDeletePost(t *testing.T) {
	prepare()
	postid := 13
	strpostid := strconv.Itoa(postid)
	post := &models.Post{
		UserID: 17,
		ID:     postid,
		Title:  "Post to test delete",
		Body:   "Post to test delete",
	}
	if err := db.SavePost(post); err != nil {
		t.Errorf("could not save post: %s", err)
	}
	if err := db.DeletePost(strpostid); err != nil {
		t.Errorf("could not delete post: %s", err)
	}
	if post, err := db.GetPost(strpostid); post != nil || err == nil {
		t.Errorf("error: post still exists")
	}
}

func TestDeletePostWithComments(t *testing.T) {
	prepare()
	postid := 13
	strpostid := strconv.Itoa(postid)
	post := &models.Post{
		UserID: 17,
		ID:     postid,
		Title:  "Post to test delete",
		Body:   "Post to test delete",
	}
	comments := []models.Comment{
		{
			Post:   post,
			PostID: postid,
			ID:     1,
			Name:   "name",
			Email:  "comment1",
			Body:   "comment1",
		},
		{
			Post:   post,
			PostID: postid,
			ID:     2,
			Name:   "name",
			Email:  "comment2",
			Body:   "comment2",
		},
	}
	if err := db.SavePost(post); err != nil {
		t.Errorf("could not save post: %s", err)
	}
	for _, comment := range comments {
		if err := db.SaveComment(&comment); err != nil {
			t.Errorf("could not save comments: %v", err)
		}
	}
	if err := db.DeletePost(strpostid); err != nil {
		t.Errorf("could not delete post: %s", err)
	}
	if post, err := db.GetPost(strpostid); post != nil || err == nil {
		t.Errorf("error: post still exists")
	}
	if comment, err := db.GetComment("1"); comment != nil || err == nil {
		t.Errorf("comment 1 still exists")
	}
	if comment, err := db.GetComment("2"); comment != nil || err == nil {
		t.Errorf("comment 2 still exists")
	}
}

func TestGetGoogleUserNotExist(t *testing.T) {
	prepare()
	_, err := db.GetGoogleUser("123")
	if err == nil {
		t.Errorf("user should not exist, expected error")
	}
}

func TestGetGoogleUser(t *testing.T) {
	id := 100
	gid := "12313123123123"
	user := &models.User{
		ID:    id,
		Email: "mail@example.com",
		Name:  "John Smith",
	}
	expected := &models.GoogleUser{
		ID:     gid,
		UserID: id,
		User:   user,
	}
	err := db.SaveGoogleUser(expected)
	if err != nil {
		t.Errorf("could not save user, expected no error, got: %v", err)
	}
	guser, err := db.GetGoogleUser(gid)
	if err != nil {
		t.Errorf("could not get saved user, expected no error, got: %v", err)
	}
	if guser.User == nil {
		t.Errorf("User field is nil, expected not nil")
	}
	if !reflect.DeepEqual(expected, guser) {
		t.Errorf("saved and received objects are not equal")
	}
}
