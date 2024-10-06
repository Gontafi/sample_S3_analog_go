package repository

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrUniqueName = errors.New("bucket name should be unique")

func createCSVWriter(filename string) (*csv.Writer, *os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(f)
	return writer, f, nil
}

func writeCSVRecord(writer *csv.Writer, record []string) {
	err := writer.Write(record)
	if err != nil {
		fmt.Println("Error writing record to CSV:", err)
	}
}

func HasBucketNameFromMetaData(name string) (bool, error) {
	file, err := os.Open("./data/buckets.csv")
	if err == os.ErrNotExist {
		_, f, err := createCSVWriter("./data/buckets.csv")
		if err != nil {
			return false, err
		}

		defer f.Close()
	} else if err != nil {
		return false, err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if row[0] == name {
			return true, nil
		}
	}

	return false, nil
}
