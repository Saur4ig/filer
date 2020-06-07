package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Saur4ig/mescon"
	"github.com/gorilla/mux"
)

const (
	portDef = "8585"

	filesUniqueDir = "files/unique"
	filesDir       = "files"

	// 5 mb
	maxFileChunkSize = 5 * (1 << 20)
)

func main() {
	port := flag.String("port", portDef, "app port")

	flag.Parse()

	listener(*port)
}

func listener(port string) {
	router := mux.NewRouter()

	router.PathPrefix("/get/").Handler(http.StripPrefix("/get/", http.FileServer(http.Dir("files/unique"))))
	router.HandleFunc("/save", ProcessHandler).Methods("PUT")
	info, err := mescon.GenAny(fmt.Sprintf("listener started\nhost - %s\nport - %s", "localhost", port))
	if err != nil {
		panic(err)
	}
	log.Println(info)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
