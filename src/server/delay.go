package server

import (
	"encoding/json"
	"fmt"
	"github.com/tuvistavie/securerandom"
)

type Delay struct {
	Body     MarshaledBytes `json:",omitempty"`
	Finished bool
	Uuid     string
	// TODO: webhook
	// webhook url: forward request to some url
}

func NewDelay() *Delay {
	uuid, err := securerandom.Uuid()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	d := &Delay{Uuid: uuid}
	// TODO: should we do this in constructor?
	Delays[uuid] = d //FIXME: possible racecondition
	return d
}

func (d *Delay) Json() []byte {
	var b []byte
	var err error
	if !d.Finished {
		fmt.Println("not finished")
		b, err = json.Marshal(d)
	} else {
		fmt.Println("finished")
		b, err = json.Marshal(d.Body)
	}
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	return b
}

func (d *Delay) Wait(resp chan []byte) {
	body := <-resp
	// when body arrives, add it to
	d.Body = body
	d.Finished = true
}
