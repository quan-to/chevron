package main

import (
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	remote_signer.AddHKPEndpoints(r)

	srv := &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: r, // Pass our instance of gorilla/mux in.
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
