package repository

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
	"triple-storage/internal/models"
)

func CreateBucket(bucketName string) error {
	err := os.MkdirAll(fmt.Sprintf("./data/%s", bucketName), os.ModePerm)
	if err != nil {
		return err
	}

	var w *csv.Writer

	f, err := os.Open("./data/buckets.csv")
	if err == os.ErrNotExist {
		w, f, err = createCSVWriter("./data/buckets.csv")
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	defer f.Close()

	writeCSVRecord(w, []string{bucketName, time.Now().String(), time.Now().String(), "active"})
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func GetBuckets() (*models.Buckets, error) {
	f, err := os.Open("./data/buckets.csv")
	if err == os.ErrNotExist {
		_, f, err = createCSVWriter("./data/buckets.csv")
		if err != nil {
			return &models.Buckets{}, err
		}
	} else if err != nil {
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

		if len(rec) > 2 {
			buckets.Bucket = append(buckets.Bucket, models.Bucket{
				Name:         rec[0],
				CreationDate: rec[1],
			})
		}
	}

	return &buckets, nil
}
