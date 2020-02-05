package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()

	for _, e := range os.Environ() {
		log.Print(e)
	}
}
