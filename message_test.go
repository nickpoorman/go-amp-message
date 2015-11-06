package ampmessage

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestAMPMessageEmpty(t *testing.T) {
	// it should be an empty message
	msg := New(nil)
	buf := msg.ToBuffer()

	if want, have := 1, len(buf); want != have {
		t.Errorf("Buffer should have a length of 1. want %#v, have %#v", want, have)
	}
}

func TestAMPMessageArgs(t *testing.T) {
	// it should add arguments
	args := []Arg{NewStringArg("foo"), NewStringArg("bar"), NewStringArg("baz")}
	msg := New(args)
	buf := msg.ToBuffer()

	if want, have := 28, len(buf); want != have {
		t.Errorf("Buffer should have a length of 28. want %#v, have %#v", want, have)
	}
}

func TestAMPMessageBuffer(t *testing.T) {
	// it should decode the message

	msg := New(nil)

	// push string
	msg.Push(NewStringArg("foo"))

	// push json
	obj := map[string]string{"foo": "bar"}
	jsonObj, err := NewJSONArg(obj)
	if err != nil {
		t.Error(err)
	}
	msg.Push(jsonObj)

	// push blob
	blob := []byte("bar")
	msg.Push(NewBlobArg(blob))

	msg2, err := NewFromBytes(msg.ToBuffer())
	if err != nil {
		t.Error(err)
	}

	if want, have := "foo", string(msg2.Args[0].Bytes()); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	var dat map[string]string
	if err := json.Unmarshal(msg2.Args[1].Bytes(), &dat); err != nil {
		t.Error(err)
	}

	if want, have := map[string]string{"foo": "bar"}, dat; !reflect.DeepEqual(want, have) {
		t.Errorf("want %#v, have %#v", want, have)
	}

	switch argType := msg2.Args[2].(type) {
	default:
		t.Errorf("want: BlobArg, have unexpected type %T\n", argType) // %T prints whatever type t has
	case *BlobArg:
		fmt.Printf("boolean %t\n", argType) // t has type bool
	}

	if want, have := "bar", string(msg2.Args[2].Bytes()); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

}
