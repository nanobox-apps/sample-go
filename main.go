package main

import (
	"github.com/blobstache/blobstache/api"
	"github.com/blobstache/blobstache/backends"
	"github.com/blobstache/blobstache/models"
	// _ "net/http/pprof"
	"runtime"
	"fmt"
	"flag"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	var be models.Storage
	
	switch selectedBackend {
	case "local":
		be = backends.NewLocalStorage(backendCredentials)
	default:
		be = backends.NewLocalStorage(backendCredentials)
	}

	err := models.Initialize(dbCredentials, be)
	if err != nil {
		panic(err)
	}

	models.CleanEmptyObjects()



	err = api.Start(port)
	fmt.Println(err)

}

var port string
var dbCredentials string
var selectedBackend string
var backendCredentials string

func init() {
	flag.StringVar(&port, "port", "8080", "Port to listen on")
	flag.StringVar(&dbCredentials, "dbCredentials", "dbname=live sslmode=disable", "Connection string for database connection")
	flag.StringVar(&selectedBackend, "backend", "local", "Backend data storage")
	flag.StringVar(&backendCredentials, "backendcredentials", "/tmp/data", "Backend data credentials")
	flag.Parse()
}