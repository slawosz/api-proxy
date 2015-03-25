package server

import (
	"encoding/json"
	"fmt"
)

type Req struct {
	Url    string
	Method string      // FIXME: not used now
	Body   interface{} // FIXME: not used now
}

type Resp struct {
	StatusCode string
	Body       MarshaledBytes
	Url        string
}

func (r *Resp) Json() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

// json marshaling in go marshals []byte to base64. So we need to implement
// json.Marshaler interface
type MarshaledBytes []byte

// this should return proper json
// TODO: watch non json respones - should be wrapped to json string?
func (m MarshaledBytes) MarshalJSON() ([]byte, error) {
	//return []byte(m), nil
	b := []byte(m)
	if len(b) == 0 {
		return []byte("null"), nil
	}
	return b, nil
}
