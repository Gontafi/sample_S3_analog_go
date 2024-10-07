package core

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"triple-storage/utils"
)

func InitDir() error {
	err := os.MkdirAll(fmt.Sprintf("./%s", utils.Directory), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fmt.Sprintf("./%s/buckets.csv", utils.Directory), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := os.Stat(fmt.Sprintf("./%s/buckets.csv", utils.Directory))
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	if fi.Size() == 0 {
		if err := writeColumnsForBucketMeta(w); err != nil {
			return err
		}
	}
	log.Println(fi.Size())
	return nil
}
