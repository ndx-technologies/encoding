package base10_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/ndx-technologies/encoding/base10"
)

func ExampleL9() {
	var v base10.L9
	v.UnmarshalText([]byte("123456789"))
	fmt.Println(v)
	// Output: 123456789
}

func TestL9(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tests := []struct {
			s string
			v uint32
		}{
			{s: "999999999", v: 999_999_999},
			{s: "123456789", v: 123_456_789},
		}
		for _, tc := range tests {
			v := base10.L9FromUint32(tc.v)
			if tc.s != v.String() {
				t.Error(tc.s, v.String())
			}
			if tc.v != v.UInt32() {
				t.Error(tc.v, v.UInt32())
			}

			x, err := base10.L9FromString(tc.s)
			if err != nil {
				t.Error(err)
			}
			if v != x {
				t.Error(v, x)
			}
		}
	})

	t.Run("error", func(t *testing.T) {
		tests := []string{
			"0123456789",
			"a",
			"",
			"1234a6789",
			"1234 6789",
		}
		for _, tc := range tests {
			_, err := base10.L9FromString(tc)
			if err == nil {
				t.Error("expected error")
			}
		}
	})
}

func FuzzL9_EncodeDecode(f *testing.F) {
	f.Add(uint32(999_999_999))
	f.Add(uint32(1))
	f.Add(uint32(0))

	f.Fuzz(func(t *testing.T, v uint32) {
		v = v % 999_999_999

		q := base10.L9FromUint32(v)
		if v != q.UInt32() {
			t.Error(v, q, q.UInt32())
		}

		g, err := base10.L9FromString(q.String())
		if err != nil {
			t.Error(err)
		}
		if q != g {
			t.Error(q, g)
		}

		b, err := q.MarshalText()
		if err != nil {
			t.Error(err)
		}

		var a base10.L9
		if !a.IsZero() {
			t.Error("expected empty")
		}
		if err := a.UnmarshalText(b); err != nil {
			t.Error(err)
		}
		if a != q {
			t.Error(a, q)
		}
	})
}

func FuzzL9_AppendBinary(f *testing.F) {
	f.Fuzz(func(t *testing.T, out []byte, v uint32) {
		v = v % 999_999_999
		q := base10.L9FromUint32(v)

		b, err := q.MarshalText()
		if err != nil {
			t.Error(err)
		}

		outBefore := make([]byte, len(out))
		copy(outBefore, out)

		out, err = q.AppendText(out)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(outBefore, out[:len(outBefore)]) {
			t.Error(outBefore, out)
		}
		if !bytes.Equal(b, out[len(outBefore):]) {
			t.Error(b, out[len(outBefore):])
		}
	})
}

func BenchmarkL9(b *testing.B) {
	v := base10.L9FromUint32(rand.Uint32())

	b.Run("string", func(b *testing.B) {
		for b.Loop() {
			v.String()
		}
	})

	b.Run("from_string", func(b *testing.B) {
		s := v.String()

		for b.Loop() {
			base10.L9FromString(s)
		}
	})
}
