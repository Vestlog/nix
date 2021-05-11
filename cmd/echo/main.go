package main

import (
	"log"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/vestlog/nix/cmd/echo/api"
	_ "github.com/vestlog/nix/cmd/echo/docs"
	"github.com/vestlog/nix/pkg/storage"
)

var (
	dsn = "storage.db"
)

// @title NIX echo API
// @version 0.1
// @description This is a sample server.
// @host localhost:8080
// @BasePath /api/v1/
func main() {
	db, err := storage.CreateGormDatabase(dsn)
	if err != nil {
		log.Fatal(err)
	}
	a := &api.EchoApi{db}
	e := echo.New()

	e.GET("/api/v1/posts", a.GetAllPosts)
	e.GET("/api/v1/posts/:id", a.GetPost)
	e.GET("/api/v1/comments", a.GetAllComments)
	e.GET("/api/v1/comments/:id", a.GetComment)
	e.GET("/api/v1/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
