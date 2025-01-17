package app

import (
	"fmt"
	"net/http"

	"triple-storage/internal/core"
	"triple-storage/internal/handlers"
)

func RunApp(port string) error {
	// bucket handlers
	http.HandleFunc("PUT /{BucketName}", handlers.PutBucketHandler)
	http.HandleFunc("GET /", handlers.GetBucketsHandler)
	http.HandleFunc("DELETE /{BucketName}", handlers.DeleteBucketHandler)

	// object handlers
	http.HandleFunc("PUT /{BucketName}/{ObjectKey}", handlers.PutObjectHandler)
	http.HandleFunc("GET /{BucketName}/{ObjectKey}", handlers.GetObjectsHandler)
	http.HandleFunc("DELETE /{BucketName}/{ObjectKey}", handlers.DeleteObjectHandler)

	if port == "" {
		port = "8080"
	}

	if err := core.InitDir(); err != nil {
		return err
	}

	fmt.Println("server started on port:", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
