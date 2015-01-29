package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Req struct {
	Url    string
	Method string
	Body   []byte
}

type Resp struct {
	StatusCode string
	Body       []byte
	Url        string
}

func main() {
	go serveMock()
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var reqs []*Req
		err := decoder.Decode(&reqs)
		if err != nil {
			panic(err)
		}
		makeRequests(reqs, w)
	})

	log.Fatal(http.ListenAndServe(":8001", nil))
}

func makeRequests(reqs []*Req, w http.ResponseWriter) http.ResponseWriter {
	responsesCh := make(chan []byte)
	wg := &sync.WaitGroup{}
	for _, req := range reqs {
		// fmt.Fprintf(w, "\n")
		// fmt.Printf("Request: %v %v \n", req.Method, req.Url)
		go func(r *Req) {
			wg.Add(1)
			fmt.Printf("Request: %v %v \n", r.Method, r.Url)
			resp, err := http.Get(r.Url)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			responsesCh <- body
			fmt.Printf("Body: %v", string(body))
			wg.Done()
		}(req)
	}

	go func() {
		wg.Wait()
		close(responsesCh)
	}()
	// responses := []*Resp{}

	// I, an idiot, was rangeing over array here instead over channel
	for resp := range responsesCh {
		fmt.Printf("-")
		//responses = append(responses, resp)
		fmt.Fprintf(w, string(resp))
	}
	fmt.Println("*")

	return w
}
