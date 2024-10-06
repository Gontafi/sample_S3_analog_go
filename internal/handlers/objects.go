package handlers

import (
	"log"
	"net/http"
	"triple-storage/internal/repository"
)

func PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("BucketName")
	objKey := r.URL.Query().Get("ObjectKey")

	ok, err := repository.HasBucketNameFromMetaData(bucketName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// some xml err
		log.Println(err)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		// w.write xml err
		return
	}

	contType := r.Header.Get("Content-Type")
	contLength := r.Header.Get("Content-Length")

	var data []byte

	_, err = r.Body.Read(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// some xml err
		log.Println(err)
		return
	}

	err = repository.AddObject(bucketName, objKey, contLength, contType, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// some xml err
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetObjectsHandler(w http.ResponseWriter, r *http.Request) {
}

func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
}
