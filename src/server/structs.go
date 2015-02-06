package server

type Req struct {
	Url    string
	Method string      // FIXME: not used now
	Body   interface{} // FIXME: not used now
}

type Resp struct {
	StatusCode string
	Body       interface{}
	Url        string
}
