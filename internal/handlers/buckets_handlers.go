package handlers

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"triple-storage/internal/core"
	"triple-storage/utils"
)

func PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]
	if !utils.IsValidBucketName(bucketName) {
		w.WriteHeader(http.StatusBadRequest)
	}

	ok, err := core.HasBucketNameFromMetaData(bucketName)
	if ok {
		w.WriteHeader(http.StatusConflict)
		// w.write xml err
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}

	err = core.CreateBucket(bucketName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	// w.Write() success
}

func GetBucketsHandler(w http.ResponseWriter, r *http.Request) {
	buckets, err := core.GetBuckets()
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
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]

	ok, err := core.HasBucketNameFromMetaData(bucketName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		// w.write xml err
		return
	}
	// check obj meta
	isEmpty, err := core.IsCSVEmpty(fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketName))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}
	if !isEmpty {
		w.WriteHeader(http.StatusConflict)

		return
	}

	err = core.DeleteBucket(bucketName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}
	// change bucket meta
	w.WriteHeader(http.StatusNoContent)
}
