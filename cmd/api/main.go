package main

import (
	"mohhefni/go-blog-app/internal/config"
)

func main() {
	filePath := ".env"
	err := config.LoadConfig(filePath)
	if err != nil {
		panic(err)
	}
}
