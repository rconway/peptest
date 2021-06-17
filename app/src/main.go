package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	// Flags
	var authMode bool
	var resourceMode bool
	var port int
	flag.BoolVar(&authMode, "auth", false, "use auth mode")
	flag.BoolVar(&resourceMode, "resource", false, "use resource mode")
	flag.IntVar(&port, "port", 80, "port to listen on")
	flag.Parse()

	router := mux.NewRouter()

	// select handler based upon the indicated mode
	var handler http.HandlerFunc
	if authMode {
		fmt.Println("auth mode enabled")
		handler = authHandler
	} else if resourceMode {
		fmt.Println("resource mode enabled")
		handler = resourceHandler
	} else {
		log.Fatal(fmt.Errorf("must specify one of -auth or -resource"))
	}

	// root
	router.PathPrefix("/").HandlerFunc(handler)

	// Start server
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("Starting to listen at %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

// Handler to provide the 'auth' response
func authHandler(w http.ResponseWriter, r *http.Request) {
	// Response body - deferred to end of function
	defer func() {
		fmt.Fprintf(w, "Endpoint = Auth Handler\n")
		dumpRequest(w, r)
		dumpRequest(os.Stdout, r)
	}()

	// Get the Authorization header
	authorizationStr := r.Header.Get("Authorization")

	// FORBIDDEN - if there is no Authorization header then treat as FORBIDDEN
	if len(authorizationStr) == 0 {
		fmt.Println("FORBIDDEN: no Authorization header")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Get the Bearer token from the Authorization header
	authParts := strings.Split(authorizationStr, " ")

	// FORBIDDEN - if there is no Bearer token
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		fmt.Println("FORBIDDEN: no Bearer token")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// UNAUTHORIZED - if the token is not 'good'
	// Also set the 'WWW-Authenticate' header
	if authParts[1] != "good" {
		fmt.Println("UNAUTHORIZED: the token is BAD")
		w.Header().Set("WWW-Authenticate", "use this ticket: xxx")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// AUTHORIZED - if all checks pass
	fmt.Println("AUTHORIZED: the token is GOOD")
	w.WriteHeader(http.StatusOK)
}

// Handler for the 'protected' resource
func resourceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoint = Resource Server\n")
	dumpRequest(w, r)
	dumpRequest(os.Stdout, r)
}

// Helper function to dump the received request to stdout
func dumpRequest(w io.Writer, r *http.Request) {
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
