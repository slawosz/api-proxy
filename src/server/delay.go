package server

import (
	"fmt"
	"github.com/tuvistavie/securerandom"
)

type Delay struct {
	Id   string
	Body []byte
	// TODO: webhook
	// webhook url: forward request to some url
}

func NewDelay() *Delay {
	// assign secure id
	uuid, err := securerandom.Uuid()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	d := &Delay{}
	Delays[uuid] = d //Fixme: possible racecondition
	return d
}

func (d *Delay) Json() []byte {
	var b []byte
	if d.Body == nil {
		// marshal some json
	} else {
		// marshal actual body
	}
	return b
}

func (d *Delay) Wait(resp chan []byte) {
	body := <-resp
	// when body arrives, add it to
	d.Body = body
}
