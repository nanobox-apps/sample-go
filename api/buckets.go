package api

import (
	"encoding/json"
	"fmt"
	"github.com/blobstache/blobstache/models"
	"net/http"
	"regexp"
)

func createBucket(rw http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")
	name := req.Header.Get("Name")
	if name == "" || userId == "" || userKey == "" {
		fmt.Println("missing params")
		rw.WriteHeader(422) // I need a bucket name
		return
	}

	buck, err := models.CreateBucket(userId, userKey, name)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(422)
		return
	}
	b, _ := json.Marshal(buck)

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}

func deleteBucket(rw http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")
	re := regexp.MustCompile("/buckets/(.*)")
	id := re.FindStringSubmatch(req.URL.Path)[1]

	err := models.DeleteBucket(userId, userKey, id)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

func getBucket(rw http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")
	re := regexp.MustCompile("/buckets/(.*)")
	id := re.FindStringSubmatch(req.URL.Path)[1]

	buck, err := models.GetBucket(userId, userKey, id)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(buck)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}

func listBuckets(rw http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")

	if userId == "" || userKey == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	bucks, err := models.ListBuckets(userId, userKey)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(bucks)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
