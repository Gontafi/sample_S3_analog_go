package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
	"triple-storage/internal/models"
	"triple-storage/utils"
)

func AddObject(bucketName, objKey, size, contType string, data io.ReadCloser) error {
	CSVpath := fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketName)

	var w *csv.Writer

	f, err := os.Open(CSVpath)
	if err == os.ErrNotExist {
		w, f, err = createCSVWriter(CSVpath)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	defer f.Close()

	writeCSVRecord(w, []string{objKey, size, contType, time.Now().String()})
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	newFile, err := os.Create(fmt.Sprintf("./%s/%s/%s", utils.Directory, bucketName, objKey))
	if err != nil {
		return err
	}

	defer newFile.Close()

	if _, err := io.Copy(newFile, data); err != nil {
		return err
	}

	return nil
}

func GetObjectMeta(bucketName, objKey string) (*models.Object, error) {
	f, err := os.Open(fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketName))
	if err == os.ErrNotExist {
		_, f, err = createCSVWriter(fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketName))
		if err != nil {
			return nil, err
		}

	} else if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(rec) > 3 {
			return &models.Object{
				Name:          rec[0],
				ContentLength: rec[1],
				ContentType:   rec[2],
			}, nil
		}
	}

	return nil, nil
}

func DeleteObject(bucketName, objKey string) error {
	err := os.Remove(fmt.Sprint("./%s/%s/%s", utils.Directory, bucketName, objKey))
	if err != nil {
		return err
	}

	err = DeleteRowInCSV(objKey, fmt.Sprint("./%s/%s/objects.csv", utils.Directory, bucketName))
	if err != nil {
		return err
	}

	return nil
}
