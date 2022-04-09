package main

import (
	"log"

	"github.com/aditya37/geofence-service/service"
)

func main() {
	srv, err := service.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	srv.Run()
}
