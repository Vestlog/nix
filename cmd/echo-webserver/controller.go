package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vestlog/nix/pkg/auth"
	"github.com/vestlog/nix/pkg/models"
	"github.com/vestlog/nix/pkg/sessions"
	"github.com/vestlog/nix/pkg/storage"
)

var (
	StoreName   = "sessionid"
	StateLength = 32
)

type Controller struct {
	DB           *storage.GormDatabase
	Store        *sessions.SessionStore
	GoogleAuth   *auth.OAuth
	FacebookAuth *auth.OAuth
	UserField    string
}

func (ctr *Controller) SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctr.Store.CreateSession(c.Response(), c.Request())
		return next(c)
	}
}

func (ctr *Controller) RestrictAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := ctr.Store.GetData(c.Request(), ctr.UserField)
		if err != nil || user == nil {
			return echo.NewHTTPError(http.StatusForbidden, "FORBIDDEN")
		}
		return next(c)
	}
}

func (ctr *Controller) GoogleLogin(c echo.Context) error {
	return ctr.Login(c, ctr.GoogleAuth)
}

func (ctr *Controller) FacebookLogin(c echo.Context) error {
	return ctr.Login(c, ctr.FacebookAuth)
}

func (ctr *Controller) Login(c echo.Context, auth *auth.OAuth) error {
	state, _ := GenerateRandomString(StateLength)
	if err := ctr.Store.SaveString(c.Response(), c.Request(), "state", state); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	url := auth.GetAuthURL(state)
	return c.Redirect(http.StatusFound, url)
}

func (ctr *Controller) GoogleCallback(c echo.Context) error {
	return ctr.Callback(c, ctr.GoogleAuth)
}

func (ctr *Controller) FacebookCallback(c echo.Context) error {
	return ctr.Login(c, ctr.FacebookAuth)
}

func (ctr *Controller) Callback(c echo.Context, auth *auth.OAuth) error {
	if err := ctr.StateOK(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	userinfo, err := auth.GetUserData(c.Request().Context(), c.QueryParam("code"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	user := &models.User{
		Email: userinfo["email"].(string),
		Name:  userinfo["name"].(string),
	}
	id, ok := userinfo["id"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "id is not a string")
	}
	guser, err := ctr.DB.GetGoogleUser(id)
	if err != nil {
		guser = &models.GoogleUser{
			ID:   id,
			User: user,
		}
		if err := ctr.DB.SaveGoogleUser(guser); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	if err := ctr.Store.SaveData(c.Response(), c.Request(), ctr.UserField, guser.User); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.Redirect(http.StatusFound, "/")
}

func (ctr *Controller) Signout(c echo.Context) error {
	if err := ctr.Store.DeleteSession(c.Response(), c.Request()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.Redirect(http.StatusFound, "/")
}

func (ctr *Controller) StateOK(c echo.Context) error {
	state, err := ctr.Store.GetString(c.Request(), "state")
	if err != nil {
		return fmt.Errorf("error: could not get session field: %w", err)
	}
	if state != c.QueryParam("state") {
		return fmt.Errorf("state is not OK")
	}
	return nil
}

func (ctr *Controller) DeletePost(c echo.Context) error {
	if err := ctr.DB.DeletePost(c.Param("postid")); err != nil {
		return fmt.Errorf("could not delete post: %w", err)
	}
	return c.Redirect(http.StatusFound, "/")
}

func (ctr *Controller) EditPostForm(c echo.Context) error {
	postid := c.Param("postid")
	post, err := ctr.DB.GetPost(postid)
	if err != nil {
		return fmt.Errorf("error getting post from database: %w", err)
	}
	data := struct {
		Action string
		Post   *models.Post
	}{
		Action: fmt.Sprintf("/%s/editpost", postid),
		Post:   post,
	}
	return c.Render(http.StatusOK, "postform", data)
}

func (ctr *Controller) EditPost(c echo.Context) error {
	postidstr := c.Param("postid")
	postid, err := strconv.Atoi(postidstr)
	if err != nil {
		return fmt.Errorf("post ID has to be an integer")
	}
	post := &models.Post{
		ID: postid,
		// Title: html.EscapeString(c.FormValue("title")),
		// Body:  html.EscapeString(c.FormValue("body")),
		Title: c.FormValue("title"),
		Body:  c.FormValue("body"),
	}
	if err := ctr.DB.UpdatePost(post); err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}
	return c.Redirect(http.StatusFound, "/"+postidstr)
}

func (ctr *Controller) CreatePostForm(c echo.Context) error {
	data := struct {
		Action string
		Post   *models.Post
	}{
		"/admin/createpost",
		new(models.Post),
	}
	return c.Render(http.StatusOK, "postform", data)
}

func (ctr *Controller) CreatePost(c echo.Context) error {
	title := c.FormValue("title")
	body := c.FormValue("body")
	if title == "" || body == "" {
		return fmt.Errorf("title or body cannot be empty")
	}
	rawuser, err := ctr.Store.GetData(c.Request(), "user")
	if err != nil {
		err = fmt.Errorf("error: could not get user data: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	user, ok := rawuser.(*models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "user is wrong")
	}
	if err := ctr.DB.SavePost(&models.Post{
		UserID: user.ID,
		Title:  title,
		Body:   body,
	}); err != nil {
		return fmt.Errorf("could not save post: %w", err)
	}
	return c.Redirect(http.StatusFound, "/")
}

func (ctr *Controller) GetAllPosts(c echo.Context) error {
	posts, err := ctr.DB.GetPosts()
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "index", posts)
}

func (ctr *Controller) GetPost(c echo.Context) error {
	id := c.Param("postid")
	post, err := ctr.DB.GetPost(id)
	if err != nil {
		return fmt.Errorf("error getting post: %w", err)
	}
	comments, err := ctr.DB.GetCommentsPostID(id)
	if err != nil {
		return fmt.Errorf("error getting comments for postid %s: %w", id, err)
	}
	rawuser, err := ctr.Store.GetData(c.Request(), ctr.UserField)
	if err != nil {
		err = fmt.Errorf("error: could not get user data: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	user, ok := rawuser.(*models.User)
	if !ok {
		user = &models.User{}
	}
	data := struct {
		Post     *models.Post
		Comments []models.Comment
		Prefix   string
		Name     string
		Email    string
	}{
		Post:     post,
		Comments: comments,
		Prefix:   "/admin",
		Name:     user.Name,
		Email:    user.Email,
	}
	return c.Render(http.StatusOK, "post", data)
}

func (ctr *Controller) CreateComment(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	body := c.FormValue("body")
	if name == "" || email == "" || body == "" {
		return fmt.Errorf("name or email or body cannot be empty")
	}
	postidstr := c.Param("postid")
	postid, err := strconv.Atoi(postidstr)
	if err != nil {
		return fmt.Errorf("post id has to be an integer")
	}
	if err := ctr.DB.SaveComment(&models.Comment{
		PostID: postid,
		Name:   name,
		Email:  email,
		Body:   body,
	}); err != nil {
		return fmt.Errorf("could not save comment for post %d : %w", postid, err)
	}
	return c.Redirect(http.StatusFound, "/"+postidstr)
}

func NotImplemented(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented yet.")
}

func CreateController(dsn string, authkey []byte) (*Controller, error) {
	db, err := storage.CreateGormDatabase(dsn)
	if err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}
	if err := db.CreateTables(); err != nil {
		return nil, fmt.Errorf("error: could not migrate database: %w", err)
	}
	store := sessions.CreateCookieStore("session", authkey)
	return &Controller{
		DB:        db,
		Store:     store,
		UserField: "user",
	}, nil
}
