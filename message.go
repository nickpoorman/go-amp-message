package ampmessage

import (
	"encoding/json"
	"fmt"

	"github.com/nickpoorman/go-amp"
)

// Arg is an abstract type that knows what type it is
type Arg interface {
	Bytes() []byte
}

// JSONArg is a JSON Arg
type JSONArg []byte

// NewJSONArg will return a new instance of JSONArg from marshalled json
func NewJSONArg(v interface{}) (*JSONArg, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	b := JSONArg(buf)
	return &b, nil
}

// Bytes returns the underlying bytes for this Arg
func (a *JSONArg) Bytes() []byte {
	return *a
}

// StringArg is a String Arg
type StringArg []byte

// NewStringArg will return a new instance of StringArg from a string
func NewStringArg(str string) *StringArg {
	b := StringArg([]byte(str))
	return &b
}

// Bytes returns the underlying bytes for this Arg
func (a *StringArg) Bytes() []byte {
	return *a
}

// BlobArg is a Blob Arg
type BlobArg []byte

// NewBlobArg simply returns the same byte array given to it
func NewBlobArg(buf []byte) *BlobArg {
	b := BlobArg(buf)
	return &b
}

// Bytes returns the underlying bytes for this Arg
func (a *BlobArg) Bytes() []byte {
	return *a
}

// Message is a high level AMP message
type Message struct {
	Args []Arg
}

// New will return a AMP Message from the provided args as bytes
func New(args []Arg) *Message {
	return &Message{
		Args: args, // we have to decode the args
	}
}

// NewFromBytes will return a AMP Message from the provided args as encoded bytes
func NewFromBytes(msg []byte) (*Message, error) {
	args, err := decode(msg)
	if err != nil {
		return nil, err
	}
	return &Message{
		Args: args, // we have to decode the args
	}, nil
}

// Push an arg onto the tail of a Message
func (m *Message) Push(arg Arg) []Arg {
	m.Args = append(m.Args, arg)
	return m.Args
}

// Pop an arg off the tail of a Message
func (m *Message) Pop() Arg {
	if len(m.Args) < 1 {
		return nil
	}
	var x Arg
	x, m.Args = m.Args[len(m.Args)-1], m.Args[:len(m.Args)-1]
	return x
}

// Shift an arg onto the head of a Message
func (m *Message) Shift() Arg {
	if len(m.Args) < 1 {
		return nil
	}
	var x Arg
	x, m.Args = m.Args[0], m.Args[1:]
	return x
}

// Unshift an arg off the head of a Message
func (m *Message) Unshift(arg Arg) []Arg {
	m.Args = append([]Arg{arg}, m.Args...)
	return m.Args
}

// Inspect the Message
func (m *Message) Inspect() string {
	return fmt.Sprintf("<Message args=%d size=%d>", len(m.Args), len(m.ToBuffer()))
}

// ToBuffer return an encoded AMP Message
func (m *Message) ToBuffer() []byte {
	return encode(m.Args)
}

// encode and pack all `args`.
func encode(args []Arg) []byte {
	byteArgs := make([][]byte, len(args))
	for i := 0; i < len(args); i++ {
		byteArgs[i] = pack(args[i])
	}
	return amp.Encode(byteArgs)
}

// pack `arg`
func pack(arg Arg) []byte {
	// undefined
	if arg == nil {
		return []byte("j:null") // JSON.stringify(null) -> "null"
	}

	switch arg.(type) {
	default: // json
		return append([]byte("j:"), arg.Bytes()...)
	case *BlobArg: // blob
		return arg.Bytes()
	case *StringArg: // string
		return append([]byte("s:"), arg.Bytes()...)
	}
}

// decode `msg` and unpack all args
func decode(msg []byte) ([]Arg, error) {
	args := amp.Decode(msg)

	var unpackedArgs []Arg
	for i := 0; i < len(args); i++ {
		unpackedArg, err := unpack(args[i])
		if err != nil {
			return nil, err
		}
		unpackedArgs = append(unpackedArgs, unpackedArg)
	}

	return unpackedArgs, nil
}

func unpack(arg []byte) (Arg, error) {
	// json
	if isJSON(arg) {
		b := JSONArg(arg[2:])
		return &b, nil
	}

	// string
	if isString(arg) {
		b := StringArg(arg[2:])
		return &b, nil
	}

	// blob
	b := BlobArg(arg)
	return &b, nil
}

// isString will return true if it's a string argument
func isString(arg []byte) bool {
	return 115 == arg[0] && 58 == arg[1]
}

// isJSON will return true if it's a JSON argument
func isJSON(arg []byte) bool {
	return 106 == arg[0] && 58 == arg[1]
}
