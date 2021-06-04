package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	authMode := flag.Bool("auth", false, "use auth mode")
	resourceMode := flag.Bool("resource", false, "use resource mode")
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	router := mux.NewRouter()

	// auth mode
	if *authMode {
		fmt.Println("auth mode enabled")
		router.PathPrefix("/auth").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			outcome, haveOutcome := r.URL.Query()["outcome"]
			if haveOutcome {
				code, err := strconv.Atoi(outcome[0])
				if err == nil {
					w.WriteHeader(code)
				}
			}

			fmt.Fprintf(w, "this is /auth\n")
			for k, v := range r.URL.Query() {
				fmt.Fprintf(w, "  %v = %v\n", k, v[0])
			}
		})
	}

	// resource mode
	if *resourceMode {
		fmt.Println("resource mode enabled")
		router.PathPrefix("/resource").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "this is /resource\n")
		})
	}

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this is root path\n")
	})

	addr := fmt.Sprintf(":%v", *port)
	fmt.Printf("Starting to listen at %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
