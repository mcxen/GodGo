package main

import (
	"log"

	"godgo/frontend"
)

func main() {
	if err := frontend.Run(); err != nil {
		log.Fatal(err)
	}
}
