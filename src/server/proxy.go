package server

import (
	"encoding/json"
	"fmt"
	"helpers"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	timeout = 1000
)

var (
	client *http.Client
	Delays map[string]*Delay
)

func init() {
	client = &http.Client{}
	Delays = make(map[string]*Delay)
}

// This functions:
// 1. Makes paraller requests
// 2. Collects responses - with delay handling
// 3. Writes response to response writer
func makeRequests(reqs []*Req, w http.ResponseWriter) http.ResponseWriter {
	t := helpers.StartTimer("*")
	responsesCh := make(chan []byte)
	wg := &sync.WaitGroup{}
	t.Check("Before loop")
	// STEP 1: Distribute requests
	for _, req := range reqs {
		resp := make(chan []byte)     // channel to collect response
		go makeRequest(req, wg, resp) // make actual request
		// handle response passing
		go func() {
			select {
			case r := <-resp: // when request finishes
				responsesCh <- r // success, distribute response to return channel
			case <-time.Tick(timeout * time.Millisecond): // when request tas long
				fmt.Println("timeout")
				// handle delay
				delay := NewDelay()
				go func() {
					delay.Wait(resp)
				}()
				// TODO: we need to send delay information to responses channel
				responsesCh <- delay.Json()
			}
			wg.Done()
		}()
	}

	// STEP 2: Wait until all responses will finish.
	// This is background thread.
	go func() {
		wg.Wait()
		close(responsesCh)
	}()

	// STEP 3: Gather all responses. It will block until all
	// request will be finished
	var responses []*Resp
	for resp := range responsesCh {
		fmt.Println(string(resp))
		responses = append(responses, &Resp{Body: MarshaledBytes(resp)})
	}
	json, err := json.Marshal(responses)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(json))

	t.Finish("Finished making requests")
	return w
}

// makes request and saves body to response
// TODO: probably we want to get channel for resp
func makeRequest(r *Req, wg *sync.WaitGroup, body chan []byte) {
	wg.Add(1)
	resp, err := client.Get(r.Url)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	body <- b
}
