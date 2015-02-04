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

	"github.com/tuvistavie/securerandom"
)

var (
	client *http.Client
	Delays map[string]Resp
)

func init() {
	client = &http.Client{}
}

type Req struct {
	Url    string
	Method string
	Body   interface{}
}

type Resp struct {
	StatusCode string
	Body       interface{}
	Url        string
}

type Timer struct {
	start int64
	prev  int64
	id    string
}

func StartTimer(id string) *Timer {
	t := &Timer{}
	t.id = id
	t.start = time.Now().UnixNano()
	t.prev = time.Now().UnixNano()
	return t
}

func (t *Timer) Check(str string) {
	now := time.Now().UnixNano()
	fmt.Printf("[+%vns] [%v] %v\n", (now - t.prev), t.id, str)
	t.prev = time.Now().UnixNano()
}

func (t *Timer) Finish(str string) {
	now := time.Now().UnixNano()
	fmt.Printf("[+%vns] [%v] %v\n", (now - t.start), t.id, str)
}

func main() {
	go serveMock()
	r := mux.NewRouter()
	r.HandleFunc("/{uuid}", callbackHandler)
	http.Handle("/callback", r)
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

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	t, ok := Delays[uuid]
	if != ok {
		http.NotFound(w,r)
		return
	}
	fmt.Fprintf(w, t.Json())
}

func makeRequests(reqs []*Req, w http.ResponseWriter) http.ResponseWriter {
	t := StartTimer("*")
	responsesCh := make(chan []byte)
	wg := &sync.WaitGroup{}
	t.Check("Before loop")
	for _, req := range reqs {
		// here we should send single channel resp
		resp := make(chan []byte)
		go makeRequest(req, wg, resp)
		select {
		// got response, lets send it to channel
		case r := <-resp:
			responsesCh <- r
		case <-time.Tick(Timeout):
			// handle timeout
			timeout := &NewTimeout{}
			go func() {
				// handle timeout,
				// should be select again
				timeout.Expect(resp)
			}()
			// TODO: send timeout to responsesCh properly
			responsesCh <- timeout.Json()
		}

	}

	t.Check("After loop")
	go func() {
		// blocks
		wg.Wait()
		close(responsesCh)
	}()

	// I, an idiot, was rangeing over array here instead over channel
	var buffer bytes.Buffer
	buffer.WriteString("[")
	t.Check("Gathering responses started")
	for resp := range responsesCh {
		buffer.Write(resp)
		buffer.WriteString(",")
	}
	t.Check("Gathering responses finishes")
	// lets remove last colon
	bb := buffer.Bytes()
	bb[len(bb)-1] = ']'
	fmt.Fprintf(w, string(bb))

	t.Finish("Finished making requests")
	return w
}

func makeRequest(r *Req, wg *sync.WaitGroup, resp <-chan []byte) {
	t := StartTimer(r.Url)
	wg.Add(1)
	// fmt.Printf("Request: %v %v \n", r.Method, r.Url)
	//resp, err := http.Get(r.Url)
	// should we implement timeout here?
	resp, err := client.Get(r.Url)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	t.Check("Request done")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	t.Check("Body read done")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	t.Check("Before pushing")
	resp <- body
	t.Check("After pushing")
	// fmt.Printf("Body: %v", string(body))
	wg.Done()
	t.Finish(fmt.Sprintf("Req finished"))
}

func handleTimeout() {
}

type Timeout struct {
	Id   string
	Body []byte
	// TODO: webhook
	// webhook url: forward request to some url
}

func NewTimeout() *Timeout {
	// assign secure id
	uuid, err := securerandom.Uuid()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// add timeout to global timeouts map
	t
	Delays[uuid] := t
}

func (t *Timeout) Json() []byte {
	var b []byte
	if t.Body == nil {
		// marshal some json
	} else {
		// marshal actual body
	}
}

func (t *Timeout) Expect(resp chan []byte) {
	// wait for body
	body <- resp
	// when body arrives, add it to
	timeout.Body = Body
}
