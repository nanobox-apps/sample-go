package api

import (
	"encoding/json"
	"fmt"
	"github.com/blobstache/blobstache/models"
	"io"
	"net/http"
	"regexp"
)

func replaceObject(rw http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")
	bucketId := req.Header.Get("Bucketid")
	re := regexp.MustCompile("/objects/(.*)")
	id := re.FindStringSubmatch(req.URL.Path)[1]
	
	obj, err := models.GetObject(userId, userKey, bucketId, id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	// create a new temporary file
	// the new file has an alias that is the id of the old object
	tmpObj, err := models.CreateObject(userId, userKey, obj.BucketID, obj.ID)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write to the tmp file
	w, err := tmpObj.WriteCloser()
	if err != nil {
		models.DeleteObject(userId, userKey, tmpObj.BucketID, tmpObj.ID)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer w.Close()
	defer req.Body.Close()

	// fill in the tmp object
	size, err := io.Copy(w, req.Body)
	if err != nil {
		if err = tmpObj.Remove(); err == nil {
			models.DeleteObject(userId, userKey, tmpObj.BucketID, tmpObj.ID)
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// move the tmp object to the existing one
	err = tmpObj.Move(obj.ID)
	if err != nil {
		if err = obj.Remove(); err == nil {
			models.DeleteObject(userId, userKey, tmpObj.BucketID, tmpObj.ID)
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set size of replaced object
	obj.Size = int64(size)
	err = models.SaveObject(obj)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = tmpObj.Remove(); err == nil {
		models.DeleteObject(userId, userKey, tmpObj.BucketID, tmpObj.ID)
	}

	f, _ := json.Marshal(obj)

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(f)
}

func createObject(rw http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")
	bucketId := req.Header.Get("Bucketid")
	alias := req.Header.Get("Alias")

	if userId == "" || userKey == "" || bucketId == "" {
		rw.WriteHeader(422)
		return
	}

	obj, err := models.CreateObject(userId, userKey, bucketId, alias)
	if err != nil {
		fmt.Println("this")
		fmt.Println(err)
		fmt.Println("end")
		rw.WriteHeader(422)
		return
	}

	w, err := obj.WriteCloser()
	if err != nil {
		fmt.Println(err)
		models.DeleteObject(userId, userKey, obj.BucketID, obj.ID)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer w.Close()
	defer req.Body.Close()

	size, err := io.Copy(w, req.Body)
	if err != nil {
		fmt.Println(err)
		if err = obj.Remove(); err == nil {
			models.DeleteObject(userId, userKey, obj.BucketID, obj.ID)
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	obj.Size = int64(size)
	err = models.SaveObject(obj)
	if err != nil {
		fmt.Println(err)
		if err = obj.Remove(); err == nil {
			models.DeleteObject(userId, userKey, obj.BucketID, obj.ID)
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, _ := json.Marshal(obj)

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(f)
}

func getObject(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("get object")
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")
	bucketId := req.Header.Get("Bucketid")
	re := regexp.MustCompile("/objects/(.*)")
	id := re.FindStringSubmatch(req.URL.Path)[1]

	obj, err := models.GetObject(userId, userKey, bucketId, id)
	if err != nil {
		rw.WriteHeader(422)
		return
	}

	rc, err := obj.ReadCloser()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer rc.Close()

	rw.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(rw, rc)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getObjectInfo(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("get object info")
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")
	bucketId := req.Header.Get("Bucketid")
	re := regexp.MustCompile("/object_info/(.*)")
	id := re.FindStringSubmatch(req.URL.Path)[1]
	
	obj, err := models.GetObject(userId, userKey, bucketId, id)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}

func deleteObject(rw http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")
	bucketId := req.Header.Get("Bucketid")
	re := regexp.MustCompile("/objects/(.*)")
	id := re.FindStringSubmatch(req.URL.Path)[1]
	obj, err := models.GetObject(userId, userKey, bucketId, id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	err =	obj.Remove()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = models.DeleteObject(userId, userKey, obj.BucketID, obj.ID)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(422)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

func listObjects(rw http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("Userid")
	userKey := req.Header.Get("Key")
	bucketId := req.Header.Get("Bucketid")
	objs, err := models.ListObjects(userId, userKey, bucketId)
	if err != nil {
		rw.WriteHeader(422)
		return
	}

	b, err := json.Marshal(objs)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
