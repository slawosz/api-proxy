package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func ServeMock() {
	r := mux.NewRouter()
	r.HandleFunc("/users/{user}/repos", usersHandler)
	http.Handle("/", r)
	http.HandleFunc("/bla", primitiveHandler)
	log.Fatal(http.ListenAndServe(":8002", nil))
}

func primitiveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bla")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixNano()
	vars := mux.Vars(r)
	user := vars["user"]

	if user == "timeout" {
		fmt.Println("using timeout")
		time.Sleep(300 * time.Millisecond)
	}

	resp := map[string]string{
		"reposCount": fmt.Sprintf("%v", len(user)),
		"reposOwner": user,
	}
	buf, _ := json.Marshal(resp)

	fmt.Fprintf(w, string(buf))
	end := time.Now().UnixNano()
	fmt.Printf("Completed in: %v ns\n", (end - start))
}
