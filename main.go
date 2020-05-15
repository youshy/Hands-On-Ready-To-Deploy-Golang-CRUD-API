package main

import (
	"log"
	"os"
)

func main() {
	a := App{}
	a.Initialize()
	a.Run(":9000")
}

func init() {
	if ok := os.Getenv("PG_USERNAME"); ok == "" {
		log.Fatalln("PG_USERNAME not specified")
	}
	if ok := os.Getenv("PG_PASSWORD"); ok == "" {
		log.Fatalln("PG_PASSWORD not specified")
	}
	if ok := os.Getenv("PG_DB_NAME"); ok == "" {
		log.Fatalln("PG_DB_NAME not specified")
	}
	if ok := os.Getenv("PG_DB_HOST"); ok == "" {
		log.Fatalln("PG_DB_HOST_ not specified")
	}
}
