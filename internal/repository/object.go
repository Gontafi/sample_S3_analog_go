package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func AddObject(bucketName, objKey, size, contType string, data []byte) error {
	CSVpath := fmt.Sprintf("./data/%s/objects.csv", bucketName)

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

	newFile, err := os.Create(fmt.Sprintf("./data/%s/%s", bucketName, objKey))
	if err != nil {
		return err
	}

	defer newFile.Close()
	_, err = newFile.Write(data)
	if err != nil {
		return err
	}

	return nil
}
