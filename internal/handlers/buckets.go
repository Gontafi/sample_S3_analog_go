package handlers

import (
	"encoding/xml"
	"log"
	"net/http"
	"triple-storage/internal/repository"
	"triple-storage/utils"
)

func PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("BucketName")

	if !utils.IsValidBucketName(bucketName) {
		w.WriteHeader(http.StatusBadRequest)
	}

	ok, err := repository.HasBucketNameFromMetaData(bucketName)
	if ok {
		w.WriteHeader(http.StatusConflict)
		// w.write xml err
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// some xml err
		log.Println(err)
		return
	}

	err = repository.CreateBucket(bucketName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// some xml err
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	// w.Write() success
}

func GetBucketsHandler(w http.ResponseWriter, r *http.Request) {
	buckets, err := repository.GetBuckets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	xmlText, err := xml.MarshalIndent(buckets, " ", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(xmlText)
}

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("BucketName")

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
	// check obj meta
	// change bucket meta
	w.WriteHeader(http.StatusNoContent)
}
