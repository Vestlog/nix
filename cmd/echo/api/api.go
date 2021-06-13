package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vestlog/nix/pkg/storage"
	"gorm.io/gorm/logger"
)

type EchoApi struct {
	DB storage.Database
}

// GetAllPosts godoc
// @Summary Get all posts
// @Produce json
// @Produce xml
// @Success 200 {array} object
// @Router /api/v1/posts [get]
func (api *EchoApi) GetAllPosts(c echo.Context) error {
	status := http.StatusOK
	var data interface{}
	var err error

	data, err = api.DB.GetPosts()
	if err != nil {
		status = http.StatusInternalServerError
		data = ErrMap(err)
	}
	return Encode(c, status, data)
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
	status := http.StatusOK
	var data interface{}

	data, err := api.DB.GetPost(id)
	if err != nil {
		status = http.StatusInternalServerError
		if err == logger.ErrRecordNotFound {
			status = http.StatusNotFound
		}
		data = ErrMap(err)
	}
	return Encode(c, status, data)
}

// GetAllComments godoc
// @Summary Get all comments
// @Produce json
// @Produce xml
// @Success 200 {array} object
// @Router /api/v1/comments [get]
func (api *EchoApi) GetAllComments(c echo.Context) error {
	status := http.StatusOK
	var data interface{}

	data, err := api.DB.GetComments()
	if err != nil {
		status = http.StatusInternalServerError
		data = ErrMap(err)
	}
	return Encode(c, status, data)
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
	status := http.StatusOK
	var data interface{}

	data, err := api.DB.GetComment(id)
	if err != nil {
		status = http.StatusInternalServerError
		if err == logger.ErrRecordNotFound {
			status = http.StatusNotFound
		}
		data = ErrMap(err)
	}
	return Encode(c, status, data)
}

func Encode(c echo.Context, status int, data interface{}) error {
	if c.Request().Header.Get("Accept") == "text/xml" {
		return c.XMLPretty(status, data, " ")
	}
	return c.JSONPretty(status, data, " ")
}

func ErrMap(err error) interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}
