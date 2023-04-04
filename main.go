package main

import (
	"fmt"
	"os"

	"test_H5B/cache"
	"test_H5B/controller"

	"github.com/labstack/echo/v4"
)

const (
	fileName = "Latin-Lipsum.txt"
	port     = 62626
)

func main() {
	b, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	memCache, err := cache.Build(b)
	if err != nil {
		panic(err)
	}
	control := controller.New(memCache)
	e := echo.New()
	searchPrefix := e.Group("/api/v0.1/search")
	searchPrefix.GET("/:searchWord", control.SearchGet)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
