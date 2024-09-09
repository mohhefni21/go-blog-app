package main

import (
	"fmt"
	"log"
	"mohhefni/go-blog-app/apps/auth"
	"mohhefni/go-blog-app/apps/post"
	"mohhefni/go-blog-app/external/database"
	"mohhefni/go-blog-app/internal/config"

	"github.com/labstack/echo/v4"
)

func main() {
	filePath := ".env"
	err := config.LoadConfig(filePath)
	if err != nil {
		panic(err)
	}

	db, err := database.Connection(config.Cfg.DBconfig)
	if err != nil {
		panic(err)
	}

	if db != nil {
		log.Println("db connected")
	}

	e := echo.New()

	// Routes
	auth.Init(e, db)
	post.Init(e, db)

	addr := fmt.Sprintf("127.0.0.1%s", config.Cfg.AppConfig.AppPort)
	fmt.Printf("starting web server at %s \n", addr)
	err = e.Start(addr)

	if err != nil {
		panic(err)
	}

}
