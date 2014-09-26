// Package main is an example of how to create an API server using
// gorilla-mux. It creates and listens to some end points
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Response is an open map to return key value pairs back to the user
type Response map[string]interface{}

// String representation of the response map
func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func main() {

	// Register a couple of routes.
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)

	log.Println("Starting up a http server on port 8080...")
	log.Println("Listening to /")
	// Send all incoming requests to mux.DefaultRouter.
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Handles requests to /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")

	fmt.Fprintln(w, "method: "+r.Method)
	fmt.Fprintln(w, "uri: "+r.RequestURI)
	fmt.Fprintln(w, "request content length: "+strconv.FormatInt(r.ContentLength, 10))

	fmt.Fprintln(w, "Headers:")
	for k, v := range r.Header {
		//values = make ([]string, 0 len(v))
		fmt.Fprint(w, "\t"+k+": ")
		if len(v) > 1 {
			fmt.Fprint(w, "\n")
			for _, hv := range v {
				fmt.Fprintln(w, "\t"+hv)
			}
		} else {
			fmt.Fprintln(w, v[0])
		}
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	fmt.Fprint(w, "Body: ")
	if len(body) > 0 {

		fmt.Fprintf(w, "\n%s", body)
	} else {
		fmt.Fprintln(w, "no body found")
	}

}
