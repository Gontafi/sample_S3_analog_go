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
	isEmpty, err := core.IsCSVEmpty(fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketName))
	if check(err, w) {
		return
	}
	if !isEmpty {
		err = core.UpdateRowInCSV(bucketName, fmt.Sprintf("./%s/buckets.csv", utils.Directory), []string{"", "", utils.CurrentTime(), "marked for delition"})
		if check(err, w) {
			return
		}

		check(core.ErrBucketIsNotEmpty, w, http.StatusConflict)
		return
	}

	err = core.DeleteBucket(bucketName)
	if check(err, w) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
