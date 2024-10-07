package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"triple-storage/internal/core"
	"triple-storage/internal/models"
	"triple-storage/utils"
)

func PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]

	if check(utils.IsValidBucketName(bucketName), w, http.StatusBadRequest) {
		return
	}

	ok, err := core.HasBucketNameFromMetaData(bucketName)
	if ok {
		check(core.ErrUniqueName, w, http.StatusConflict)
		return
	}
	if check(err, w) {
		return
	}

	err = core.CreateBucket(bucketName)
	if check(err, w) {
		return
	}

	w.WriteHeader(http.StatusOK)
	// w.Write() success
}

func GetBucketsHandler(w http.ResponseWriter, r *http.Request) {
	buckets, err := core.GetBuckets()
	if check(err, w) {
		return
	}
	xmlText, err := xml.MarshalIndent(models.ListAllMyBucketsResult{
		XMLName: xml.Name{},
		Buckets: *buckets,
	}, " ", " ")
	if check(err, w) {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(xmlText)
}

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]

	ok, err := core.HasBucketNameFromMetaData(bucketName)
	if check(err, w) {
		return
	}
	if !ok {
		check(utils.ErrNotFound, w, http.StatusNotFound)
		return
	}
	// check obj meta
	isEmpty, err := core.IsCSVEmpty(fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketName))
	if check(err, w) {
		return
	}
	if !isEmpty {
		check(core.ErrBucketIsNotEmpty, w, http.StatusConflict)
		return
	}

	err = core.DeleteBucket(bucketName)
	if check(err, w) {
		return
	}
	// change bucket meta
	w.WriteHeader(http.StatusNoContent)
}
