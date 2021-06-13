package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/vestlog/nix/pkg/models"
	mock "github.com/vestlog/nix/pkg/storage/mock_storage"
	"gorm.io/gorm/logger"
)

func TestGetPost(t *testing.T) {
	post := &models.Post{
		UserID: 1,
		ID:     1,
		Title:  "title",
		Body:   "text",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockDatabase(ctrl)
	m.EXPECT().GetPost(gomock.Eq("1")).Return(post, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	api := &EchoApi{
		DB: m,
	}
	if err := api.GetPost(c); err != nil {
		t.Error(err)
	}
	r := &models.Post{}
	json.NewDecoder(rec.Body).Decode(r)
	if !reflect.DeepEqual(post, r) {
		t.Errorf("expected %v, got %v", post, r)
	}
}

func TestGetPostErrNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockDatabase(ctrl)
	m.
		EXPECT().
		GetPost(gomock.Eq("1")).
		Return(nil, logger.ErrRecordNotFound)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	api := &EchoApi{
		DB: m,
	}
	if err := api.GetPost(c); err != nil {
		t.Error(err)
	}

	if rec.Result().StatusCode != http.StatusNotFound {
		t.Errorf("got %v, expected %v",
			rec.Result().StatusCode,
			http.StatusNotFound)
	}
}

func TestGetAllPosts(t *testing.T) {
	posts := make([]models.Post, 10)
	for i := 0; i < len(posts); i++ {
		posts[i] = models.Post{
			UserID: 10,
			ID:     i,
			Title:  fmt.Sprintf("title %d", i),
			Body:   "text",
		}
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockDatabase(ctrl)
	m.EXPECT().GetPosts().Return(posts, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/")
	api := &EchoApi{
		DB: m,
	}
	if err := api.GetAllPosts(c); err != nil {
		t.Error(err)
	}
	var r []models.Post
	if err := json.NewDecoder(rec.Body).Decode(&r); err != nil {
		t.Errorf("could not decode json: %v", err)
	}
	if !reflect.DeepEqual(posts, r) {
		t.Errorf("expected %v, got %v", posts, r)
	}
}

func TestGetAllPostsErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockDatabase(ctrl)
	m.
		EXPECT().
		GetPosts().
		Return(nil, errors.New("some error"))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/")
	api := &EchoApi{
		DB: m,
	}
	if err := api.GetAllPosts(c); err != nil {
		t.Error(err)
	}

	if rec.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("got %v, expected %v",
			rec.Result().StatusCode,
			http.StatusNotFound)
	}
}

func TestGetComment(t *testing.T) {
	comment := &models.Comment{
		Post:   nil,
		PostID: 1,
		ID:     1,
		Name:   "name",
		Email:  "mail@example.com",
		Body:   "text",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockDatabase(ctrl)
	m.EXPECT().GetComment(gomock.Eq("1")).Return(comment, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	api := &EchoApi{
		DB: m,
	}
	if err := api.GetComment(c); err != nil {
		t.Error(err)
	}
	r := &models.Comment{}
	json.NewDecoder(rec.Body).Decode(r)
	if !reflect.DeepEqual(comment, r) {
		t.Errorf("expected %v, got %v", comment, r)
	}
}

func TestGetCommentErrNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockDatabase(ctrl)
	m.
		EXPECT().
		GetComment(gomock.Eq("1")).
		Return(nil, logger.ErrRecordNotFound)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	api := &EchoApi{
		DB: m,
	}
	if err := api.GetComment(c); err != nil {
		t.Error(err)
	}

	if rec.Result().StatusCode != http.StatusNotFound {
		t.Errorf("got %v, expected %v",
			rec.Result().StatusCode,
			http.StatusNotFound)
	}
}

func TestGetAllComments(t *testing.T) {
	comments := make([]models.Comment, 10)
	for i := 0; i < len(comments); i++ {
		comments[i] = models.Comment{
			Post:   nil,
			PostID: 1,
			ID:     i,
			Name:   "name",
			Email:  "mail@example.com",
			Body:   fmt.Sprintf("text %d", i),
		}
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockDatabase(ctrl)
	m.EXPECT().GetComments().Return(comments, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/")
	api := &EchoApi{
		DB: m,
	}
	if err := api.GetAllComments(c); err != nil {
		t.Error(err)
	}
	var result []models.Comment
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Errorf("could not decode json: %v", err)
	}
	if !reflect.DeepEqual(comments, result) {
		t.Errorf("expected %v, got %v", comments, result)
	}
}

func TestGetAllCommentsErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockDatabase(ctrl)
	m.
		EXPECT().
		GetComments().
		Return(nil, errors.New("some error"))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/")
	api := &EchoApi{
		DB: m,
	}
	if err := api.GetAllComments(c); err != nil {
		t.Error(err)
	}

	if rec.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("got %v, expected %v",
			rec.Result().StatusCode,
			http.StatusNotFound)
	}
}
