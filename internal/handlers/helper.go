package handlers

import (
	"encoding/xml"
	"log"
	"net/http"
	"triple-storage/internal/models"
)

func check(err error, w http.ResponseWriter, code ...int) bool {
	if err != nil {
		log.Println(err)
		if len(code) == 0 {
			code = append(code, 500)
		}

		xmlText, err := xml.MarshalIndent(models.Error{
			Code:    code[0],
			Message: err.Error(),
		}, " ", " ")
		if err != nil {
			log.Println(err)
		}

		w.WriteHeader(code[0])
		w.Write(xmlText)

		return true
	}

	return false
}
