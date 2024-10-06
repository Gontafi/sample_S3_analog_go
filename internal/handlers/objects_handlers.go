package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"triple-storage/internal/core"
)

func PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]
	objKey := parts[2]

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

	contType := r.Header.Get("Content-Type")
	contLength := r.Header.Get("Content-Length")

	size, err := strconv.ParseInt(contLength, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}

	data := make([]byte, 0, size)

	_, err = r.Body.Read(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}

	err = core.AddObject(bucketName, objKey, contLength, contType, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetObjectsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	bucketName := parts[1]
	objKey := parts[2]

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

	ok, err = core.HasObjkeyInMeta(bucketName, objKey)
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

	obj, err := core.GetObjectMeta(bucketName, objKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}

	if obj == nil {
		w.WriteHeader(http.StatusInternalServerError)
		// w.write xml err
		return
	}

	file, err := os.Open(fmt.Sprintf("./data/%s/%s", bucketName, obj.Name))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// w.write xml err
		log.Println(err)
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

	ok, err = core.HasObjkeyInMeta(bucketName, objKey)
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

	err = core.DeleteObject(bucketName, objKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
