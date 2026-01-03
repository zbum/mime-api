package main

import (
	"log"
	"net/http"
	"os"

	"github.com/zbum/mantyboot/http/mux"
	"github.com/zbum/mantyboot/http/mux/middleware"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	mantyMux := mux.NewMantyMux()
	mantyMux.AddMiddleware(middleware.AccessLogger(logger))
	mantyMux.HandleFunc("POST /v1/display-part", DisplayPart)
	mantyMux.HandleFunc("GET /health", Health)
	log.Fatal(http.ListenAndServe(":8080", mantyMux))
}
