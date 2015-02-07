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
	Delays map[string]*Delay
	client *http.Client
)

func init() {
	client = &http.Client{}
}

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
			case <-time.Tick(timeout * time.Millisecond): // when request thaks long
				fmt.Println("timeout")
				// handle delay
				delay := NewDelay()
				go func() {
					delay.Wait(resp)
				}()
				// TODO: we need to send delay information to responses channel
				responsesCh <- delay.Json()
			}
		}()
	}

	t.Check("After loop")
	// STEP 2: Wait until all responses will finish.
	// This is background thread.
	go func() {
		fmt.Println("before wait")
		wg.Wait()
		fmt.Println("after wait")
		close(responsesCh)
	}()

	// STEP 3: Gather all responses. It will block until all
	// request will be finished
	var responses []*Resp
	for resp := range responsesCh {
		responses = append(responses, &Resp{Body: resp})
	}
	json, err := json.Marshal(responses)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(json))

	t.Finish("Finished making requests")
	return w
}

func makeRequest(r *Req, wg *sync.WaitGroup, body chan []byte) {
	t := helpers.StartTimer(r.Url)
	wg.Add(1)
	defer wg.Done()
	resp, err := client.Get(r.Url)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	t.Check("Request done")
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	t.Check("Body read done")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	t.Check("Before pushing")
	body <- b
	t.Check("After pushing")
	t.Finish(fmt.Sprintf("Req finished"))
}
