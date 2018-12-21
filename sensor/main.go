package main

import (
	"log"
)

func main() {
	s := server{}
	if err := s.run(); err != nil {
		log.Fatal(err)
	}
}
