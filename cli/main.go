package main

import (
	// "net/http"
	"flag"
	"fmt"
)

func main() {
	if (len(flag.Args()) < 1) {
		fmt.Println("What should I do?")
		return
	}

	switch flag.Arg(0) {
  case "list-users", "list_users", "listusers":
    listUsers()
  case "create-user", "create_user", "createuser":
    createUser()
  case "delete-user", "delete_user", "deleteuser":
    deleteUser(flag.Arg(1))
  case "list-buckets", "list_buckets", "listbuckets":
    listBuckets()
  case "show-bucket", "show_bucket", "showbucket":
    showBucket(flag.Arg(1))
  case "create-bucket", "create_bucket", "createbucket":
    createBucket()
  case "delete-bucket", "delete_bucket", "deletebucket":
    deleteBucket(flag.Arg(1))
  case "list-objects", "list_objects", "listobjects":
    listObjects()
  case "object-size", "object_size", "objectsize":
    showObjectSize(flag.Arg(1))
  case "get-object", "get_object", "getobject":
    getObject(flag.Arg(1))
  case "create-object", "create_object", "createobject":
    createObject()
  case "delete-object", "delete_object", "deleteobject":
    deleteObject(flag.Arg(1))

  default:
    fmt.Printf("I dont know what to do with %s\n", flag.Arg(0))
  }
}