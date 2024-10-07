package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"triple-storage/internal/models"
	"triple-storage/utils"
)

func CreateBucket(bucketName string) error {
	bucketPath := fmt.Sprintf("./%s/%s", utils.Directory, bucketName)
	err := os.MkdirAll(bucketPath, os.ModePerm)
	if err != nil {
		return err
	}

	csvFilePath := fmt.Sprintf("./%s/buckets.csv", utils.Directory)
	var w *csv.Writer
	var f *os.File

	f, err = os.OpenFile(csvFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	w = csv.NewWriter(f)

	fi, err := os.Stat(csvFilePath)
	if err != nil {
		return err
	}

	if fi.Size() == 0 {
		if err := writeColumnsForBucketMeta(w); err != nil {
			return err
		}
	}

	record := []string{bucketName, utils.CurrentTime(), utils.CurrentTime(), "active"}
	if err := writeCSVRecord(w, record); err != nil {
		return err
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	objCSVFile, err := os.OpenFile(fmt.Sprintf("%s/objects.csv", bucketPath), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return err
	}

	defer objCSVFile.Close()

	fi, err = os.Stat(fmt.Sprintf("%s/objects.csv", bucketPath))
	if err != nil {
		return err
	}

	w = csv.NewWriter(objCSVFile)

	if fi.Size() == 0 {
		if err := writeColumnsForObjMeta(w); err != nil {
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func GetBuckets() (*models.Buckets, error) {
	f, err := os.Open(fmt.Sprintf("./%s/buckets.csv", utils.Directory))
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	buckets := models.Buckets{
		Bucket: []models.Bucket{},
	}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(rec) > 3 {
			buckets.Bucket = append(buckets.Bucket, models.Bucket{
				Name:             rec[0],
				CreationDate:     rec[1],
				LastModifiedTime: rec[2],
				Status:           rec[3],
			})
		}
	}
	if len(buckets.Bucket) > 1 {
		buckets.Bucket = buckets.Bucket[1:]
	}

	return &buckets, nil
}

func DeleteBucket(bucketName string) error {
	err := DeleteRowInCSV(bucketName, fmt.Sprintf("./%s/buckets.csv", utils.Directory))
	if err != nil {
		return err
	}
	return os.RemoveAll(fmt.Sprintf("./%s/%s", utils.Directory, bucketName))
}
