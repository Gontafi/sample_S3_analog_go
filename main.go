package main

import (
	"flag"
	"fmt"
	"log"
	"triple-storage/internal/app"
	"triple-storage/utils"
)

func main() {
	port := flag.String("port", "8080", "")
	dir := flag.String("dir", "data", "")

	flag.Usage = func() {
		fmt.Println(utils.Help)
	}

	flag.Parse()

	utils.Directory = *dir

	err := app.RunApp(*port)
	if err != nil {
		log.Fatal(err)
	}
}
