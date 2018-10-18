// +build ignore

package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addrF := flag.String("addr", "127.0.0.1:8080", "Address to listen")
	dirF := flag.String("dir", ".", "Directory to serve")
	flag.Parse()

	log.Printf("Starting server on http://%s/ ...", *addrF)

	http.Handle("/", http.FileServer(http.Dir(*dirF)))
	http.ListenAndServe(*addrF, nil)
}
