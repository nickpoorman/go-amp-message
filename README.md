# amp-message

  High level [AMP](https://github.com/nickpoorman/go-amp) `Message` implementation for manipulating, encoding and decoding AMP messages in Go.

  [Godoc](https://godoc.org/github.com/nickpoorman/go-amp-message)

## Installation

```
$ go get github.com/nickpoorman/go-amp-message
```

## Example

  Encoding a message:

```go
msg := ampmessage.New(nil)

fmt.Printf("<Bytes: %x>\n", msg.ToBytes())
// => <Bytes: 10>

msg.Push(ampmessage.NewStringArg("foo"))
msg.PushString("bar") // convenience method for pushing a new StringArg
msg.PushString("baz")
fmt.Printf("<Bytes: %x>\n", msg.ToBytes())
// => <Bytes: 1300000005733a666f6f00000005733a62617200000005733a62617a>

// json
jsonObj, _ := ampmessage.NewJSONArg(map[string]string{"foo": "bar"})
msg.Push(jsonObj)
fmt.Printf("<Bytes: %x>\n", msg.ToBytes())
// => <Bytes: 1400000005733a666f6f00000005733a62617200000005733a62617a0000000f6a3a7b22666f6f223a22626172227d>

// convenience method for pushing a new JSONArg
msg.PushJSON(map[string]string{"ping": "pong"})
fmt.Printf("<Bytes: %x>\n", msg.ToBytes())
// => <Bytes: 1500000005733a666f6f00000005733a62617200000005733a62617a0000000f6a3a7b22666f6f223a22626172227d000000116a3a7b2270696e67223a22706f6e67227d>

blob := []byte("beep")
msg.Push(ampmessage.NewBlobArg(blob))
msg.PushBlob([]byte("boop")) // convenience method for pushing a new BlobArg
fmt.Printf("<Bytes: %x>\n", msg.ToBytes())
// => <Bytes: 1700000005733a666f6f00000005733a62617200000005733a62617a0000000f6a3a7b22666f6f223a22626172227d000000116a3a7b2270696e67223a22706f6e67 ... >
```

  Decoding a message:

```go
msg := ampmessage.New(nil)

msg.PushString("foo")
msg.PushJSON(map[string]string{"hello": "world"})
msg.PushBlob([]byte("hello"))

other, _ := ampmessage.NewFromBytes(msg.ToBytes())

fmt.Printf("%s\n", other.Shift())
// => &foo
fmt.Printf("%s\n", other.Shift())
// => &{"hello":"world"}
fmt.Printf("%s\n", other.Shift())
// => &hello
```

## API

[Godoc](https://godoc.org/github.com/nickpoorman/go-amp-message)

### `ampmessage.New(nil)`

  Initialize an empty message.

### `ampmessage.NewFromBytes(bytes)`

  Decode the `buffer` AMP message to populate the `Message`.

### `ampmessage.New(args)`

  Initialize a message populated with `args`.

# License

  MIT
