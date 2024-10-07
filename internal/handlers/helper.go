package handlers

import (
	"log"
	"net/http"
)

func check(err error, w http.ResponseWriter) bool {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println(err)

		return true
	}

	return false
}
