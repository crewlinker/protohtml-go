package phtml

import (
	"github.com/go-playground/form/v4"
)

// PHTML is embeddded in any generated handlers.
type PHTML struct {
	dec ValuesDecoder
	enc ValuesEncoder
}

// New initializes the handler set.
func New() *PHTML {
	dec := form.NewDecoder()
	dec.SetTagName("json")
	enc := form.NewEncoder()
	enc.SetTagName("json")

	return &PHTML{
		dec: dec,
		enc: enc,
	}
}

// ValuesDecoder returns the values decoder.
func (bh PHTML) ValuesDecoder() ValuesDecoder {
	return bh.dec
}

// ValuesEncoder returns the values encoder.
func (bh PHTML) ValuesEncoder() ValuesEncoder {
	return bh.enc
}
