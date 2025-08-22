package base32_test

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"testing"

	base32x "github.com/ndx-technologies/encoding/base32"
)

func ExampleUint32_MarshalText() {
	u := base32x.Uint32(42)
	b, _ := u.MarshalText()
	fmt.Print(string(b))
	// Output: aaaaabk
}

func ExampleUint32_UnmarshalText() {
	var u base32x.Uint32
	u.UnmarshalText([]byte("aaaaabk"))
	fmt.Print(uint32(u))
	// Output: 42
}

func ExampleUint32_UnmarshalText_mixed_case() {
	var u base32x.Uint32
	u.UnmarshalText([]byte("AaaaaBk"))
	fmt.Print(uint32(u))
	// Output: 42
}

func ExampleUint32_UnmarshalText_error() {
	var u base32x.Uint32
	fmt.Print(u.UnmarshalText([]byte("aaaaaaaaaabk")))
	// Output: invalid length
}

func ExampleUint32_String() {
	u := base32x.Uint32(42)
	fmt.Print(u.String())
	// Output: aaaaabk
}

func ExampleUint32FromString() {
	u, _ := base32x.Uint32FromString("aaaaabk")
	fmt.Print(uint32(u))
	// Output: 42
}

func FuzzUint32EncodeDecode(f *testing.F) {
	f.Add(uint32(0))
	f.Add(uint32(11))
	f.Add(uint32(32123))
	f.Add(uint32(222))

	f.Fuzz(func(t *testing.T, v uint32) {
		b, err := base32x.Uint32(v).MarshalText()
		if err != nil {
			t.Error(err)
		}

		var x base32x.Uint32
		if err := x.UnmarshalText(b); err != nil || uint32(x) != v {
			t.Error(x, uint32(x), v, err)
		}
	})
}

func FuzzUint32AppendText(f *testing.F) {
	f.Add([]byte(nil), uint32(11))
	f.Add([]byte{1, 2, 3}, uint32(11))

	f.Fuzz(func(t *testing.T, out []byte, v uint32) {
		outBefore := make([]byte, len(out))
		copy(outBefore, out)

		out, err := base32x.Uint32(v).AppendText(out)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(outBefore, out[:len(outBefore)]) {
			t.Error(outBefore, out)
		}

		bb, err := base32x.Uint32(v).MarshalText()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(bb, out[len(outBefore):]) {
			t.Error(bb, out[len(outBefore):])
		}
	})
}

func BenchmarkUint32_AppendText(b *testing.B) {
	v := base32x.Uint32(rand.Uint32())
	out := make([]byte, 0, 7)

	for b.Loop() {
		v.AppendText(out)
	}
}

func BenchmarkUint32_MarshalText(b *testing.B) {
	v := base32x.Uint32(rand.Uint32())

	for b.Loop() {
		v.MarshalText()
	}
}

func BenchmarkUint32_MarshalText_standard(b *testing.B) {
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)

	v := make([]byte, 7)
	i := rand.Uint32()

	for b.Loop() {
		enc.Encode(v, binary.BigEndian.AppendUint32(nil, i))
	}
}

func BenchmarkUint32_UnmarshalText(b *testing.B) {
	var v base32x.Uint32
	s, _ := base32x.Uint32(rand.Uint32()).MarshalText()

	for b.Loop() {
		v.UnmarshalText(s)
	}
}

func BenchmarkUint32_UnmarshalText_standard(b *testing.B) {
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)

	v := make([]byte, 7)
	enc.Encode(v, binary.BigEndian.AppendUint32(nil, rand.Uint32()))

	for b.Loop() {
		var g [4]byte
		enc.Decode(g[:], v)
		binary.BigEndian.Uint32(g[:])
	}
}
