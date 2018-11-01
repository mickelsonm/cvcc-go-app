package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ninnemana/cvcc-go-app/router"
)

var (
	mux router.Router
)

func main() {

	var err error
	mux, err = router.NewBasic()
	if err != nil {
		log.Fatalf("failed to create HTTP router: %v", err)
	}

	http.HandleFunc("/put", mux.Put)
	http.HandleFunc("/add", mux.Add)
	http.HandleFunc("/", mux.Index)

	port := ":8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("fell out of serving HTTP traffic: %v", err)
	}
}
