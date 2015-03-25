package server

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNotFinishedDelaySerialization(t *testing.T) {
	d := &Delay{Uuid: "foobar"}

	expected := `{"Uuid":"foobar"}`
	b, err := json.Marshal(d)

	if err != nil {
		fmt.Println(err)
		//panic(err)
	}

	res := string(b)
	if res != expected {
		t.Errorf("%v should be equal %v", res, expected)
	}
}
