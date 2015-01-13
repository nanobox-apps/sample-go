package main

import (
	"github.com/blobstache/blobstache/api"
	"github.com/blobstache/blobstache/backends"
	"github.com/blobstache/blobstache/models"
	"runtime"
	"fmt"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := models.Initialize("dbname=live sslmode=disable", backends.NewLocalStorage("./data"))
	if err != nil {
		panic(err)
	}

	models.CleanEmptyObjects()

	err = api.Start("8080")
	fmt.Println(err)
}
