package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"triple-storage/internal/models"
	"triple-storage/utils"
)

func AddObject(bucketName, objKey, size, contType string, data io.ReadCloser) error {
	CSVpath := fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketName)

	var w *csv.Writer

	f, err := os.OpenFile(CSVpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	w = csv.NewWriter(f)

	fi, err := os.Stat(CSVpath)
	if err != nil {
		return err
	}

	if fi.Size() == 0 {
		if err := writeColumnsForObjMeta(w); err != nil {
			return err
		}
	}

	ok, err := HasObjkeyInMeta(bucketName, objKey)
	if err != nil {
		return err
	}

	if ok {
		err = DeleteRowInCSV(objKey, fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketName))
		if err != nil {
			return err
		}
	}
	writeCSVRecord(w, []string{objKey, size, contType, utils.CurrentTime()})
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
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	first := true
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if first {
			first = false
			continue
		}
		if len(rec) > 3 && rec[0] == objKey {
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
	err := os.Remove(fmt.Sprintf("./%s/%s/%s", utils.Directory, bucketName, objKey))
	if err != nil {
		return err
	}

	err = DeleteRowInCSV(objKey, fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketName))
	if err != nil {
		return err
	}

	return nil
}
