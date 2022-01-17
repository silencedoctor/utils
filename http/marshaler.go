package http

import "io"

// MarshalHandler defines a conversion between byte sequence and data interface.
// 暂不支持， 后续考虑通过content type 实现多种marshal方式
type MarshalHandler interface {
	// Marshal marshals "v" into byte sequence.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal unMarshals "data" into "v".
	// "v" must be a pointer value.
	Unmarshal(data []byte, v interface{}) error
	// NewDecoder returns a Decoder which reads byte sequence from "r".
	NewDecoder(r io.Reader) Decoder
	// NewEncoder returns an Encoder which writes bytes sequence into "w".
	NewEncoder(w io.Writer) Encoder
}

// Decoder decodes a byte sequence
type Decoder interface {
	Decode(v interface{}) error
}

// Encoder encodes data interface into byte sequence.
type Encoder interface {
	Encode(v interface{}) error
}
