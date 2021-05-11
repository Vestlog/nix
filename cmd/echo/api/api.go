package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vestlog/nix/pkg/storage"
)

type EchoApi struct {
	DB *storage.GormDatabase
}

// GetAllPosts godoc
// @Summary Get all posts
// @Produce json
// @Produce xml
// @Success 200 {array} object
// @Router /api/v1/posts [get]
func (api *EchoApi) GetAllPosts(c echo.Context) error {
	data, err := api.DB.GetPosts()
	if err != nil {
		return fmt.Errorf("error getting posts: %w", err)
	}
	return Encode(c, data)
}

// GetPost godoc
// @Summary Get post from ID
// @Description Get a single post for a given ID
// @Produce json
// @Produce xml
// @Param id path int true "post id"
// @Success 200 {object} object
// @Router /api/v1/posts/{id} [get]
func (api *EchoApi) GetPost(c echo.Context) error {
	id := c.Param("id")
	data, err := api.DB.GetPost(id)
	if err != nil {
		return fmt.Errorf("error getting post %v: %w", id, err)
	}
	return Encode(c, data)
}

// GetAllComments godoc
// @Summary Get all comments
// @Produce json
// @Produce xml
// @Success 200 {array} object
// @Router /api/v1/comments [get]
func (api *EchoApi) GetAllComments(c echo.Context) error {
	data, err := api.DB.GetComments()
	if err != nil {
		return fmt.Errorf("error getting posts: %w", err)
	}
	return Encode(c, data)
}

// GetComment godoc
// @Summary Get comment from ID
// @Description Get a single comment for a given ID
// @Produce json
// @Produce xml
// @Param id path int true "comment id"
// @Success 200 {object} object
// @Router /api/v1/comments/{id} [get]
func (api *EchoApi) GetComment(c echo.Context) error {
	id := c.Param("id")
	data, err := api.DB.GetComment(id)
	if err != nil {
		return fmt.Errorf("error getting post %v: %w", id, err)
	}
	return Encode(c, data)
}

func Encode(c echo.Context, data interface{}) error {
	if c.Request().Header.Get("Accept") == "text/xml" {
		return c.XMLPretty(http.StatusOK, data, " ")
	}
	return c.JSONPretty(http.StatusOK, data, " ")
}
