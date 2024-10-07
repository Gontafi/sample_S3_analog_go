package core

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"triple-storage/utils"
)

var (
	ErrUniqueName          = errors.New("bucket name should be unique")
	ErrBucketIsNotEmpty    = errors.New("current bucket is not empty")
	ErrInvalidLengthCSVRow = errors.New("invalid length csv rows")
)

func writeCSVRecord(writer *csv.Writer, record []string) error {
	return writer.Write(record)
}

func writeColumnsForBucketMeta(writer *csv.Writer) error {
	return writer.Write([]string{"Name", "CreationTime", "LastModifiedTime", "Status"})
}

func writeColumnsForObjMeta(writer *csv.Writer) error {
	return writer.Write([]string{"ObjectKey", "Size", "ContentType", "LastModified"})
}

func IsCSVEmpty(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	cnt, err := utils.LineCounter(file)
	if err != nil {
		log.Println(err)
	}

	return cnt < 2, nil
}

func searchKeyInCSV(f *os.File, key string) (bool, error) {
	reader := csv.NewReader(f)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if row[0] == key {
			return true, nil
		}
	}

	return false, nil
}

func HasBucketNameFromMetaData(name string) (bool, error) {
	bucketDir := fmt.Sprintf("./%s", utils.Directory)
	bucketFilePath := fmt.Sprintf("%s/buckets.csv", bucketDir)

	if _, err := os.Stat(bucketDir); os.IsNotExist(err) {
		err := os.MkdirAll(bucketDir, 0o755)
		if err != nil {
			return false, err
		}
	}

	file, err := os.OpenFile(bucketFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return false, err
	}
	defer file.Close()

	return searchKeyInCSV(file, name)
}

func HasObjkeyInMeta(bucketname, objKey string) (bool, error) {
	path := fmt.Sprintf("./%s/%s/objects.csv", utils.Directory, bucketname)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return false, err
	}

	defer file.Close()

	return searchKeyInCSV(file, objKey)
}

func DeleteRowInCSV(name, csvPath string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}

	defer f.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Split(scanner.Text(), ",")[0] != name {
			_, err := buf.Write(scanner.Bytes())
			if err != nil {
				return err
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	err = os.WriteFile(csvPath, buf.Bytes(), 0o666)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRowInCSV(name, csvPath string, record []string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}

	defer f.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rec := strings.Split(scanner.Text(), ",")
		if rec[0] == name {
			if len(rec) != len(record) {
				return ErrInvalidLengthCSVRow
			}
			for i, v := range record {
				if v != "" {
					rec[i] = record[i]
				}
			}

			_, err := buf.WriteString(strings.Join(rec, ","))
			if err != nil {
				return err
			}
		} else {
			_, err := buf.Write(scanner.Bytes())
			if err != nil {
				return err
			}
		}
		_, err = buf.WriteString("\n")
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	err = os.WriteFile(csvPath, buf.Bytes(), 0o666)
	if err != nil {
		return err
	}

	return nil
}
