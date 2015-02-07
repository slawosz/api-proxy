package server

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

// json marshaling in go marshals []byte to base64. So we need to implement
// json.Marshaler interface
type MarshaledBytes []byte

func (m MarshaledBytes) MarshalJSON() ([]byte, error) {
	return []byte(m), nil
}
