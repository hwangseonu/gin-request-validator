package main

import (
	"github.com/hwangseonu/gin-restful-example/basic/server"
	"log"
)

func main() {
	s := server.CreateServer()
	if err := s.Run(":5000"); err != nil {
		log.Fatal(err)
	}
}