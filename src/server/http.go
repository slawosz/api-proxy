// http server
package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"
)

// serves http server, and blocks
func Start() {
	fmt.Println("Server started on 8001")
	// callbacks
	r := mux.NewRouter()
	r.HandleFunc("{uuid}", delayHandler)
	http.Handle("/callback", r)
	// proxy
	http.HandleFunc("/call", proxyHandler)
	// stack
	http.HandleFunc("/stack", stacktraceHandler)

	log.Fatal(http.ListenAndServe(":8001", nil))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixNano()
	decoder := json.NewDecoder(r.Body)
	var reqs []*Req
	err := decoder.Decode(&reqs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Prepared in: %v ns\n", (time.Now().UnixNano() - start))
	makeRequests(reqs, w) // here is core functionality of the project
	end := time.Now().UnixNano()
	fmt.Printf("Completed in: %v ns\n", (end - start))
}

// handles delayed requests
func delayHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calling delay")
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	t, ok := Delays[uuid]
	if !ok {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, string(t.Json()))
}

func stacktraceHandler(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 10)
	finished := false
	var copied int
	for !finished {
		copied = runtime.Stack(buf, true)
		if copied < len(buf) {
			finished = true
		} else {
			buf = make([]byte, (len(buf)+1)*2)
		}
	}

	io.WriteString(w, string(buf))
}
