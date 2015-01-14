package api

import (
	"fmt"
	"github.com/blobstache/blobstache/models"
	"net/http"
	"regexp"
)

// Start
func Start(port string) error {
	routes := registerRoutes()

	// blocking...
	http.Handle("/", routes)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return err
	}

	return nil
}

// handleRequest
func handleRequest(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println(`
Request:
--------------------------------------------------------------------------------
%+v

`, req)

		//
		fn(rw, req)
		fmt.Println(`
Response:
--------------------------------------------------------------------------------
%+v

`, rw)
	}
}

func adminAccess(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		userId := req.Header.Get("Userid")
		userKey := req.Header.Get("Key")

		if userId == "" || userKey == "" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		// get a user and return it
		user, _ := models.GetUser(userId)
		if user == nil || user.Key != userKey || user.Admin == false {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		fn(rw, req)
	}
}


func userId(req *http.Request) string {
	return req.Header.Get("Userid")
}

func userKey(req *http.Request) string {
	return req.Header.Get("Key")	
}

func bucketId(req *http.Request) (id string) {
	id = req.Header.Get("Bucketname")
	if id == "" {
		id = req.Header.Get("Bucketid")
	}
	if id == "" {
		re := regexp.MustCompile("/buckets/(.*)")
		res := re.FindStringSubmatch(req.URL.Path)
		if len(res) == 2 {
			id = res[1]
		}
	}
	return
}

func objectId(req *http.Request) (id string) {
	id = req.Header.Get("Objectalias")
	if id == "" {
		re := regexp.MustCompile("/objects/(.*)")
		res := re.FindStringSubmatch(req.URL.Path)
		if len(res) == 2 {
			id = res[1]
		}
	}
	return
}