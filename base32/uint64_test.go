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

func ExampleUint64_MarshalText() {
	u := base32x.Uint64(42)
	b, _ := u.MarshalText()
	fmt.Print(string(b))
	// Output: aaaaaaaaaaabk
}

func ExampleUint64_UnmarshalText() {
	var u base32x.Uint64
	u.UnmarshalText([]byte("aaaaaaaaaaabk"))
	fmt.Print(uint64(u))
	// Output: 42
}

func ExampleUint64_UnmarshalText_error() {
	var u base32x.Uint64
	fmt.Print(u.UnmarshalText([]byte("aaaaaaaaaabk")))
	// Output: invalid length
}

func ExampleUint64_String() {
	u := base32x.Uint64(42)
	fmt.Println(u.String())
	// Output: aaaaaaaaaaabk
}

func ExampleUint64FromString() {
	u, _ := base32x.Uint64FromString("aaaaaaaaaaabk")
	fmt.Println(uint64(u))
	// Output: 42
}

func FuzzUint64EncodeDecode(f *testing.F) {
	f.Add(uint64(0))
	f.Add(uint64(11))
	f.Add(uint64(32123))
	f.Add(uint64(222))

	f.Fuzz(func(t *testing.T, v uint64) {
		b, err := base32x.Uint64(v).MarshalText()
		if err != nil {
			t.Error(err)
		}

		var x base32x.Uint64
		if err := x.UnmarshalText(b); err != nil || uint64(x) != v {
			t.Error(x, uint64(x), v, err)
		}
	})
}

func FuzzUint64AppendText(f *testing.F) {
	f.Add([]byte(nil), uint64(11))
	f.Add([]byte{1, 2, 3}, uint64(11))

	f.Fuzz(func(t *testing.T, out []byte, v uint64) {
		outBefore := make([]byte, len(out))
		copy(outBefore, out)

		out, err := base32x.Uint64(v).AppendText(out)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(outBefore, out[:len(outBefore)]) {
			t.Error(outBefore, out)
		}

		bb, err := base32x.Uint64(v).MarshalText()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(bb, out[len(outBefore):]) {
			t.Error(bb, out[len(outBefore):])
		}
	})
}

func BenchmarkUint64_AppendText(b *testing.B) {
	v := base32x.Uint64(rand.Uint64())
	out := make([]byte, 0, 13)

	for b.Loop() {
		v.AppendText(out)
	}
}

func BenchmarkUint64_MarshalText(b *testing.B) {
	v := base32x.Uint64(rand.Uint64())

	for b.Loop() {
		v.MarshalText()
	}
}

func BenchmarkUint64_MarshalText_standard(b *testing.B) {
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)

	v := make([]byte, 13)
	i := rand.Uint64()

	for b.Loop() {
		enc.Encode(v, binary.BigEndian.AppendUint64(nil, i))
	}
}

func BenchmarkUint64_UnmarshalText(b *testing.B) {
	var v base32x.Uint64
	s, _ := base32x.Uint64(rand.Uint64()).MarshalText()

	for b.Loop() {
		v.UnmarshalText(s)
	}
}

func BenchmarkUint64_UnmarshalText_standard(b *testing.B) {
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)

	v := make([]byte, 13)
	enc.Encode(v, binary.BigEndian.AppendUint64(nil, rand.Uint64()))

	for b.Loop() {
		var g [13]byte
		enc.Decode(g[:], v)
		binary.BigEndian.Uint64(g[:])
	}
}
