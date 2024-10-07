package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"triple-storage/internal/core"
	"triple-storage/utils"
)

func PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]
	objKey := parts[2]

	if !utils.IsValidBucketName(objKey) {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	ok, err := core.HasBucketNameFromMetaData(bucketName)
	if check(err, w) {
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		// w.write xml err
		return
	}

	contType := r.Header.Get("Content-Type")
	contLength := r.Header.Get("Content-Length")

	size, err := strconv.ParseInt(contLength, 10, 64)
	if check(err, w) {
		return
	}

	data := make([]byte, 0, size)

	_, err = r.Body.Read(data)
	if check(err, w) {
		return
	}

	err = core.AddObject(bucketName, objKey, contLength, contType, r.Body)
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
		w.WriteHeader(http.StatusNotFound)
		// w.write xml err
		return
	}

	ok, err = core.HasObjkeyInMeta(bucketName, objKey)
	if check(err, w) {
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		// w.write xml err
		return
	}

	obj, err := core.GetObjectMeta(bucketName, objKey)
	if check(err, w) {
		return
	}

	if obj == nil {
		w.WriteHeader(http.StatusInternalServerError)
		// w.write xml err
		return
	}

	file, err := os.Open(fmt.Sprintf("./data/%s/%s", bucketName, obj.Name))
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
		w.WriteHeader(http.StatusNotFound)
		// w.write xml err
		return
	}

	ok, err = core.HasObjkeyInMeta(bucketName, objKey)
	if check(err, w) {
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		// w.write xml err
		return
	}

	err = core.DeleteObject(bucketName, objKey)
	if check(err, w) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
