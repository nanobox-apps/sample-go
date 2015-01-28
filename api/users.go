package api

import (
	"encoding/json"
	"github.com/blobstache/blobstache/models"
	"net/http"
	"github.com/jcelliott/lumber"
)

func createUser(rw http.ResponseWriter, req *http.Request) {
	newUser, err := models.CreateUser()
	if err != nil {
		lumber.Error("Create User: Create :%s",err.Error())
		rw.WriteHeader(422)
		return
	}

	b, _ := json.Marshal(newUser)
	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}

func deleteUser(rw http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get(":id")
	err := models.DeleteUser(id)
	if err != nil {
		lumber.Error("Delete User: Delete :%s",err.Error())
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

func listUsers(rw http.ResponseWriter, req *http.Request) {
	users, err := models.ListUsers()
	if err != nil {
		lumber.Error("List User: Get :%s",err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(users)
	if err != nil {
		lumber.Error("List User: Json Marshel :%s",err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
