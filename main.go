package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/auth").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outcome, ok := r.URL.Query()["outcome"]
		if ok {
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

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this is root path\n")
	})

	http.ListenAndServe(":8080", router)
}
