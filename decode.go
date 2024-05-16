package jsonseal

import (
	"encoding/json"
	"io"
)

// Decoder is a drop-in replacement for the standard library json.Decoder
type Decoder struct {
	d *json.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	d := json.NewDecoder(r)
	return &Decoder{
		d: d,
	}
}

func (dec *Decoder) UseNumber() { dec.d.UseNumber() }

func (dec *Decoder) DisallowUnknownFields() { dec.d.DisallowUnknownFields() }

func (dec *Decoder) Decode(v any) error { return dec.d.Decode(v) }

func (dec *Decoder) Buffered() io.Reader { return dec.d.Buffered() }

func (dec *Decoder) InputOffset() int64 { return dec.d.InputOffset() }

func (dec *Decoder) More() bool { return dec.d.More() }

func (dec *Decoder) Token() (json.Token, error) { return dec.d.Token() }
