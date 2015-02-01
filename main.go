package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
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

type Timer struct {
	start int64
	prev  int64
}

func (t *Timer) Start() {
	t.start = time.Now().UnixNano()
	t.prev = time.Now().UnixNano()
}

func main() {
	go serveMock()
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UnixNano()
		decoder := json.NewDecoder(r.Body)
		var reqs []*Req
		err := decoder.Decode(&reqs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Prepared in: %v ns\n", (time.Now().UnixNano() - start))
		makeRequests(reqs, w)
		end := time.Now().UnixNano()
		fmt.Printf("Completed in: %v ns\n", (end - start))
	})

	log.Fatal(http.ListenAndServe(":8001", nil))
}

func makeRequests(reqs []*Req, w http.ResponseWriter) http.ResponseWriter {
	responsesCh := make(chan []byte)
	wg := &sync.WaitGroup{}
	for _, req := range reqs {
		// here we should send single channel singleResp
		go makeRequest(req, wg, responsesCh)
		/*
			select {
			// got response, lets send it to channel
			when resp := <- singleResp:
			  responsesCh <- resp
			when <- time.Tick(Timeout):
			  // handle timeout
				timeout := &NewTimeout{}
			  go func() {
					// handle timeout,
					// should be select again
					handleTimeout(timeout, singleResp)
				}
				responsesCh <- timeout
			}

		*/
	}

	go func() {
		// blocks
		wg.Wait()
		close(responsesCh)
	}()

	// I, an idiot, was rangeing over array here instead over channel
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for resp := range responsesCh {
		buffer.Write(resp)
		buffer.WriteString(",")
	}
	// lets remove last colon
	bb := buffer.Bytes()
	bb[len(bb)-1] = ']'
	fmt.Fprintf(w, string(bb))

	return w
}

func makeRequest(r *Req, wg *sync.WaitGroup, responsesCh chan []byte) {
	wg.Add(1)
	// fmt.Printf("Request: %v %v \n", r.Method, r.Url)
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
	// fmt.Printf("Body: %v", string(body))
	wg.Done()
}
