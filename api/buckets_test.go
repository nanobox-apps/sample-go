package api

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
)

func TestCreateBucket(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Userid", adminUser().ID)
	req.Header.Add("Key", adminUser().Key)
	req.Header.Add("Name", "test")

	w := httptest.NewRecorder()
	createBucket(w, req)
	if w.Code != 201 {
		t.Error("Bucket was not created:", w)
	}

}

func TestListBucket(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Userid", adminUser().ID)
	req.Header.Add("Key", adminUser().Key)

	w := httptest.NewRecorder()
	listBuckets(w, req)
	if w.Code != 200 {
		t.Error("Bucket was not created:", w)
	}
}


func TestGetBucket(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/buckets/test", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Userid", adminUser().ID)
	req.Header.Add("Key", adminUser().Key)

	w := httptest.NewRecorder()
	getBucket(w, req)
	if w.Code != 200 {
		t.Error("Bucket was not created:", w)
	}

}

func TestDeleteBucket(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://example.com/buckets/test", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Userid", adminUser().ID)
	req.Header.Add("Key", adminUser().Key)

	w := httptest.NewRecorder()
	getBucket(w, req)
	fmt.Println(w.Code)
	if w.Code != 200 {
		t.Error("Bucket not deleted:", w)
	}
}


