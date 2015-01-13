package api

import "github.com/gorilla/pat"

func registerRoutes() *pat.Router {
	router := pat.New()

	router.Post("/objects", handleRequest(createObject))
	router.Get("/object_info/{id}", handleRequest(getObjectInfo))
	router.Get("/objects/{id}", handleRequest(getObject))
	router.Put("/objects/{id}", handleRequest(replaceObject))
	router.Get("/objects", handleRequest(listObjects))
	router.Delete("/objects/{id}", handleRequest(deleteObject))

	router.Post("/buckets", handleRequest(createBucket))
	router.Get("/buckets/{id}", handleRequest(getBucket))
	router.Get("/buckets", handleRequest(listBuckets))
	router.Delete("/buckets/{id}", handleRequest(deleteBucket))

	// admin only
	router.Get("/users", handleRequest(adminAccess(listUsers)))
	router.Post("/users", handleRequest(adminAccess(createUser)))
	router.Delete("/users/{id}", handleRequest(adminAccess(deleteUser)))

	return router
}
