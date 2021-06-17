package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

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

	var handler http.HandlerFunc

	// auth mode
	if authMode {
		fmt.Println("auth mode enabled")
		// router.PathPrefix("/auth").HandlerFunc(authHandler)
		handler = authHandler
	} else if resourceMode {
		fmt.Println("resource mode enabled")
		// router.PathPrefix("/resource").HandlerFunc(resourceHandler)
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

// Handler for root path
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is root path\n")
	dumpRequest(w, r)
	dumpRequest(os.Stdout, r)
}

// Handler to provide the 'auth' response
func authHandler(w http.ResponseWriter, r *http.Request) {
	// Response body - deferred to end of function
	defer func() {
		fmt.Fprintf(w, "this is /auth\n")
		dumpRequest(w, r)
		dumpRequest(os.Stdout, r)
	}()

	// Get the original URI requested, upon which the decision is to be made
	origUriStr := r.Header.Get("X-Original-Uri")
	if len(origUriStr) == 0 {
		fmt.Println("ERROR: cannot get origUri")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse the URI string to a URL structure
	origUri, err := url.Parse(origUriStr)
	if err != nil {
		fmt.Printf("ERROR: Cannot parse origUriStr = %v\n", origUriStr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the required 'outcome' from the original request
	fmt.Printf("Query = %v\n", origUri.Query())
	outcome, haveOutcome := origUri.Query()["outcome"]
	if haveOutcome {
		code, err := strconv.Atoi(outcome[0])
		if err == nil {
			// If the code is 401 (Unauthorized) then return the expected header key
			// nginx will pass this back through to the client
			if code == 401 {
				w.Header().Set("WWW-Authenticate", "use this ticket: xxx")
			}
			w.WriteHeader(code)
		}
	}
}

// Handler for the 'protected' resource
func resourceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is /resource\n")
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
