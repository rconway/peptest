package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
			// for k, v := range r.URL.Query() {
			// 	fmt.Fprintf(w, "  %v = %v\n", k, v[0])
			// }
			dumpRequest(w, r)
		})
	}

	// resource mode
	if *resourceMode {
		fmt.Println("resource mode enabled")
		router.PathPrefix("/resource").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "this is /resource\n")
			dumpRequest(w, r)
		})
	}

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this is root path\n")
	})

	addr := fmt.Sprintf(":%v", *port)
	fmt.Printf("Starting to listen at %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func dumpRequest(w http.ResponseWriter, r *http.Request) {
	// Host
	fmt.Fprintln(w, "Host:", r.Host)
	// URL
	fmt.Fprintln(w, "URL:", r.URL)
	// Method
	fmt.Fprintln(w, "Method:", r.Method)
	// Headers
	fmt.Fprintln(w, "Headers:")
	for headerKey, headerVal := range r.Header {
		fmt.Fprintf(w, "  %v: %v\n", headerKey, headerVal)
	}
	// Query params...
	fmt.Fprintln(w, "Params:")
	for paramKey, paramVal := range r.URL.Query() {
		fmt.Fprint(w, "  ", paramKey, ":")
		for _, paramValItem := range paramVal {
			fmt.Fprint(w, " ", paramValItem)
		}
		fmt.Fprintln(w)
	}
	// Body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR reading request body:", err)
	} else {
		fmt.Fprintln(w, "Body:", string(data))
	}
}
