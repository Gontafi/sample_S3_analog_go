package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"triple-storage/internal/core"
	"triple-storage/utils"
)

func PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]
	objKey := parts[2]

	if check(utils.IsValidBucketName(bucketName), w, http.StatusBadRequest) {
		return
	}

	ok, err := core.HasBucketNameFromMetaData(bucketName)
	if check(err, w) {
		return
	}
	if !ok {
		check(utils.ErrNotFound, w, http.StatusNotFound)
		return
	}

	contType := r.Header.Get("Content-Type")
	contLength := r.Header.Get("Content-Length")

	err = core.AddObject(bucketName, objKey, contLength, contType, r.Body)
	if check(err, w) {
		return
	}

	err = core.UpdateRowInCSV(bucketName, fmt.Sprintf("./%s/buckets.csv", utils.Directory), []string{"", "", utils.CurrentTime(), ""})
	if check(err, w) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetObjectsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]
	objKey := parts[2]

	ok, err := core.HasBucketNameFromMetaData(bucketName)
	if check(err, w) {
		return
	}
	if !ok {
		check(utils.ErrNotFound, w, http.StatusNotFound)
		return
	}

	ok, err = core.HasObjkeyInMeta(bucketName, objKey)
	if check(err, w) {
		return
	}
	if !ok {
		check(utils.ErrNotFound, w, http.StatusNotFound)
		return
	}

	obj, err := core.GetObjectMeta(bucketName, objKey)
	if check(err, w) {
		return
	}

	if obj == nil {
		check(errors.New("object is nil"), w)
		return
	}

	file, err := os.Open(fmt.Sprintf("./data/%s/%s", bucketName, objKey))
	if check(err, w) {
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", obj.ContentType)
	w.Header().Set("Content-Length", obj.ContentLength)

	http.ServeContent(w, r, obj.Name, time.Now(), file)
}

func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]
	objKey := parts[2]

	ok, err := core.HasBucketNameFromMetaData(bucketName)
	if check(err, w) {
		return
	}
	if !ok {
		check(utils.ErrNotFound, w, http.StatusNotFound)
		return
	}

	ok, err = core.HasObjkeyInMeta(bucketName, objKey)
	if check(err, w) {
		return
	}
	if !ok {
		check(utils.ErrNotFound, w, http.StatusNotFound)
		return
	}

	err = core.DeleteObject(bucketName, objKey)
	if check(err, w) {
		return
	}

	err = core.UpdateRowInCSV(bucketName, fmt.Sprintf("./%s/buckets.csv", utils.Directory), []string{"", "", utils.CurrentTime(), ""})
	if check(err, w) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
