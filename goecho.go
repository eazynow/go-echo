// Package main is an example of how to create an API server using
// gorilla-mux. It creates and listens to some end points
package goecho

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ResponseParsed struct {
	RequestTime   string
	Method        string
	URI           string
	ContentLength int64
	Body          string
	RemoteAddr    string
	Header        map[string][]string
}

func (resp *ResponseParsed) parseRequest(r *http.Request) {
	resp.RequestTime = time.Now().Local().String()
	resp.Method = r.Method
	resp.URI = r.RequestURI
	resp.RemoteAddr = r.RemoteAddr
	resp.Header = r.Header
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		resp.Body = "error reading body: " + err.Error()
	}

	if len(body) > 0 {
		resp.Body = string(body[:r.ContentLength])
	} else {
		resp.Body = ""
	}
}

func (rp *ResponseParsed) AsWriter(w http.ResponseWriter) {
	fmt.Fprintf(w, "request time: %s\n", rp.RequestTime)
	fmt.Fprintln(w, "method: "+rp.Method)
	fmt.Fprintln(w, "uri: "+rp.URI)
	fmt.Fprintln(w, "request content length: "+strconv.FormatInt(rp.ContentLength, 10))

	fmt.Fprintln(w, "Headers:")
	for k, v := range rp.Header {
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
	fmt.Fprintln(w, "Body:")
	fmt.Fprintln(w, rp.Body)
}

func (rp *ResponseParsed) getAppStoreKey(c appengine.Context) *datastore.key {
	return datastore.NewKey(c, "IncomingRequest", 0, nil)
}

func (rp *ResponseParsed) save(r *http.Request) {
	c := appengine.NewContext(r)

	key := datastore.NewIncompleteKey(c, "IncomingRequest", rp.getAppStoreKey(c))
	_, err := datastore.put(c, key, &rp)
	if err != nil {
		log.Println(err.Error())
	}
}

func init() {

	log.Println("Starting up a http server on port 8080...")
	log.Println("Listening to /")

	http.HandleFunc("/", rootHandler)
}

// Handles requests to /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")

	resp := ResponseParsed()

	resp.parseRequest(r)

	resp.AsWriter(w)
}
