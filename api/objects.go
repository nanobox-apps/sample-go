package api

import (
	"encoding/json"
	"github.com/blobstache/blobstache/models"
	"io"
	"github.com/jcelliott/lumber"
	"net/http"
	"strconv"
)

func replaceObject(rw http.ResponseWriter, req *http.Request) {
	obj, err := models.GetObject(userId(req), userKey(req), bucketId(req), objectId(req))
	if err != nil {
		lumber.Error("Replace Object: Get Existing :%s",err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	// create a new temporary file
	// the new file has an alias that is the id of the old object
	tmpObj, err := models.CreateObject(userId(req), userKey(req), obj.BucketID, obj.ID)
	if err != nil {
		lumber.Error("Replace Object: New Object :%s",err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write to the tmp file
	w, err := tmpObj.WriteCloser()
	if err != nil {
		lumber.Error("Replace Object: Write to tmp :%s",err.Error())
		models.DeleteObject(userId(req), userKey(req), tmpObj.BucketID, tmpObj.ID)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer w.Close()
	defer req.Body.Close()

	// fill in the tmp object
	size, err := io.Copy(w, req.Body)
	if err != nil {
		lumber.Error("Replace Object: Copy to tmp :%s",err.Error())
		if err = tmpObj.Remove(); err == nil {
			models.DeleteObject(userId(req), userKey(req), tmpObj.BucketID, tmpObj.ID)
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// move the tmp object to the existing one
	err = tmpObj.Move(obj.ID)
	if err != nil {
		lumber.Error("Replace Object: Move Tmp :%s",err.Error())
		if err = obj.Remove(); err == nil {
			models.DeleteObject(userId(req), userKey(req), tmpObj.BucketID, tmpObj.ID)
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	// I have no need for the db record anymore
	models.DeleteObject(userId(req), userKey(req), tmpObj.BucketID, tmpObj.ID)

	// set size of replaced object
	obj.Size = int64(size)
	err = models.SaveObject(obj)
	if err != nil {
		lumber.Error("Replace Object: Save Existing :%s",err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}


	f, _ := json.Marshal(obj)

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(f)
}

func createObject(rw http.ResponseWriter, req *http.Request) {
	_, err := models.GetObject(userId(req), userKey(req), bucketId(req), objectId(req))
	// If the object already exists replace it
	if err == nil {
		replaceObject(rw, req)
		return
	}

	obj, err := models.CreateObject(userId(req), userKey(req), bucketId(req), objectId(req))
	if err != nil {
		lumber.Error("Create Object: Create :%s",err.Error())		
		rw.WriteHeader(422)
		return
	}

	w, err := obj.WriteCloser()
	if err != nil {
		lumber.Error("Create Object: Get writecloser :%s",err.Error())		
		models.DeleteObject(userId(req), userKey(req), obj.BucketID, obj.ID)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer w.Close()
	defer req.Body.Close()

	size, err := io.Copy(w, req.Body)
	if err != nil {
		lumber.Error("Create Object: Copy :%s",err.Error())		
		if err = obj.Remove(); err == nil {
			models.DeleteObject(userId(req), userKey(req), obj.BucketID, obj.ID)
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	obj.Size = int64(size)
	err = models.SaveObject(obj)
	if err != nil {
		lumber.Error("Create Object: Save :%s",err.Error())
		if err = obj.Remove(); err == nil {
			models.DeleteObject(userId(req), userKey(req), obj.BucketID, obj.ID)
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
	obj, err := models.GetObject(userId(req), userKey(req), bucketId(req), objectId(req))
	if err != nil {
		lumber.Error("Get Object: Get :%s",err.Error())		
		rw.WriteHeader(422)
		return
	}
	if obj.Size == 0 {
		lumber.Info("object size is 0", obj.Size)
		rw.WriteHeader(422)
		rw.Write([]byte("incomplete file"))
		return
	}

	rc, err := obj.ReadCloser()
	if err != nil {
		lumber.Error("Get Object: Get ReadCloser :%s",err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer rc.Close()

	rw.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(rw, rc)
	if err != nil {
		lumber.Error("Get Object: Copy :%s",err.Error())		
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getObjectInfo(rw http.ResponseWriter, req *http.Request) {
	obj, err := models.GetObject(userId(req), userKey(req), bucketId(req), objectId(req))
	if err != nil {
		lumber.Error("Get Object Info: Get :%s",err.Error())		
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Object-Alias", obj.Alias)
	rw.Header().Set("Object-Size", strconv.FormatInt(obj.Size, 10))
}

func deleteObject(rw http.ResponseWriter, req *http.Request) {
	obj, err := models.GetObject(userId(req), userKey(req), bucketId(req), objectId(req))
	if err != nil {
		lumber.Error("Delete Object: Get :%s",err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	err = obj.Remove()
	if err != nil {
		lumber.Error("Delete Object: Remove :%s",err.Error())
		if obj.Size != 0 {
			// if the object size is 0 dont worry about a failed remove
			// chances are the object didnt have any data in it.
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = models.DeleteObject(userId(req), userKey(req), obj.BucketID, obj.ID)
	if err != nil {
		lumber.Error("Delete Object: Delete :%s",err.Error())
		rw.WriteHeader(422)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

func listObjects(rw http.ResponseWriter, req *http.Request) {
	objs, err := models.ListObjects(userId(req), userKey(req), bucketId(req))
	if err != nil {
		lumber.Error("List Object: Get :%s",err.Error())
		rw.WriteHeader(422)
		return
	}

	b, err := json.Marshal(objs)
	if err != nil {
		lumber.Error("List Object: Json Marshal :%s",err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
