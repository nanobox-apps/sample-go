package api

import (
	"fmt"
	"github.com/blobstache/blobstache/models"
	"net/http"
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
