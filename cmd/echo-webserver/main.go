package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vestlog/nix/pkg/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var (
	conffilepath = flag.String("conf", "conf.json", "path to configuration file")
)

func main() {
	flag.Parse()
	LoadConfig(*conffilepath)

	goauth := &oauth2.Config{
		ClientID:     GlobalConfig.GoogleOAuth.ClientID,
		ClientSecret: GlobalConfig.GoogleOAuth.ClientSecret,
		Endpoint:     endpoints.Google,
		RedirectURL:  "http://localhost:8080/google/callback",
		Scopes:       []string{"openid", "email", "profile"},
	}
	fbauth := &oauth2.Config{
		ClientID:     GlobalConfig.FacebookOAuth.ClientID,
		ClientSecret: GlobalConfig.FacebookOAuth.ClientSecret,
		Endpoint:     endpoints.Facebook,
		RedirectURL:  "http://localhost:8080/facebook/callback",
		Scopes:       []string{"public_profile", "email"},
	}

	ctr, err := CreateController(
		GlobalConfig.DSN,
		[]byte(GlobalConfig.SessionsKey),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctr.GoogleAuth = &auth.OAuth{
		Conf:    goauth,
		InfoURL: auth.GoogleInfoURL,
	}
	ctr.FacebookAuth = &auth.OAuth{
		Conf:    fbauth,
		InfoURL: auth.FacebookInfoURL,
	}
	e := echo.New()
	// e.Debug = true
	e.Renderer = CreateTemplate()

	e.Use(middleware.Logger())
	e.Use(ctr.SessionMiddleware)

	e.GET("/", ctr.GetAllPosts)
	e.GET("/:postid", ctr.GetPost)

	restricted := e.Group("/admin")
	restricted.Use(ctr.RestrictAccess)
	restricted.GET("/:postid/editpost", ctr.EditPostForm)
	restricted.POST("/:postid/editpost", ctr.EditPost)
	restricted.GET("/createpost", ctr.CreatePostForm)
	restricted.POST("/createpost", ctr.CreatePost)
	restricted.POST("/:postid/addcomment", ctr.CreateComment)
	restricted.GET("/:postid/deletepost", ctr.DeletePost)
	restricted.GET("/signout", ctr.Signout)

	gauth := e.Group("/google")
	gauth.GET("/login", ctr.GoogleLogin)
	gauth.GET("/callback", ctr.GoogleCallback)

	fauth := e.Group("/facebook")
	fauth.GET("/login", ctr.FacebookLogin)
	fauth.GET("/callback", ctr.FacebookCallback)

	e.GET("/favicon.ico", NotImplemented)
	e.Static("/static", "static")

	e.Logger.Fatal(e.Start(":" + GlobalConfig.Port))
}
