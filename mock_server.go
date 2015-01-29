package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func serveMock() {
	r := mux.NewRouter()
	r.HandleFunc("/users/{user}/repos", usersHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8002", nil))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(".")
	vars := mux.Vars(r)
	user := vars["user"]

	if user == "timeout" {
		time.Sleep(150 * time.Millisecond)
	}

	resp := map[string]string{
		"reposCount": fmt.Sprintf("%v", len(user)),
		"reposOwner": user,
	}
	buf, _ := json.Marshal(resp)

	fmt.Fprintf(w, string(buf))
}
