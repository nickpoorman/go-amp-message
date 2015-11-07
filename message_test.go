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
	buf := msg.ToBytes()

	if want, have := 1, len(buf); want != have {
		t.Errorf("Buffer should have a length of 1. want %#v, have %#v", want, have)
	}
}

func TestAMPMessageArgs(t *testing.T) {
	// it should add arguments
	args := []Arg{NewStringArg("foo"), NewStringArg("bar"), NewStringArg("baz")}
	msg := New(args)
	buf := msg.ToBytes()

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

	msg2, err := NewFromBytes(msg.ToBytes())
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

func TestExampleEncode(t *testing.T) {
	msg := New(nil)

	fmt.Printf("<Bytes: %x>\n", msg.ToBytes())

	// string
	msg.Push(NewStringArg("foo"))
	msg.PushString("bar") // convenience method for pushing a new StringArg
	msg.PushString("baz")

	fmt.Printf("<Bytes: %x>\n", msg.ToBytes())

	// json
	jsonObj, err := NewJSONArg(map[string]string{"foo": "bar"})
	if err != nil {
		t.Error(err)
	}
	msg.Push(jsonObj)
	fmt.Printf("<Bytes: %x>\n", msg.ToBytes())

	// convenience method for pushing a new JSONArg
	if _, err := msg.PushJSON(map[string]string{"ping": "pong"}); err != nil {
		t.Error(err)
	}
	fmt.Printf("<Bytes: %x>\n", msg.ToBytes())

	// blob
	blob := []byte("beep")
	msg.Push(NewBlobArg(blob))
	msg.PushBlob([]byte("boop")) // convenience method for pushing a new BlobArg
	fmt.Printf("<Bytes: %x>\n", msg.ToBytes())
}

func TestExampleDecode(t *testing.T) {
	msg := New(nil)

	msg.PushString("foo")
	msg.PushJSON(map[string]string{"hello": "world"})
	msg.PushBlob([]byte("hello"))

	other, err := NewFromBytes(msg.ToBytes())
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%s\n", other.Shift())
	// => &foo
	fmt.Printf("%s\n", other.Shift())
	// => &&{"hello":"world"}
	fmt.Printf("%s\n", other.Shift())
	// => &hello
}
