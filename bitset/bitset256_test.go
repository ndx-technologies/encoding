package bitset_test

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/ndx-technologies/encoding/bitset"
)

func ExampleBitSet256() {
	var a bitset.BitSet256
	a.Set(11, true)
	a.Set(1000, true)

	fmt.Print(a.IsZero(), a.Get(11), a.Get(64), a.Get(256))
	// Output: false true false false
}

func ExampleBitSet256_IsZero() {
	var empty bitset.BitSet256
	fmt.Print(empty.IsZero())
	// Output: true
}

func ExampleBitSet256_Union() {
	var a, b bitset.BitSet256

	a.Set(0, true)
	b.Set(1, true)

	a.Union(b)

	fmt.Print(a.Get(0), a.Get(1), a.Get(2), a.Get(3))
	// Output: true true false false
}

func TestBitSet256_encodingError(t *testing.T) {
	var a bitset.BitSet256
	if err := a.UnmarshalBinary([]byte{0, 1, 2, 3, 4, 5}); err == nil {
		t.Error("expected error")
	}
}

func FuzzBitSet256_AppendBinary(f *testing.F) {
	f.Add([]byte(nil), []byte{1, 2, 3, 4, 5, 6, 7, 8})
	f.Add([]byte{1, 2, 3}, []byte{1, 2, 3, 4, 5, 6, 7, 8})

	f.Fuzz(func(t *testing.T, out, data []byte) {
		if len(data) < 32 {
			data = append(data, make([]byte, 32-len(data))...)
		} else if len(data) > 32 {
			data = data[:32]
		}

		var a bitset.BitSet256

		if err := a.UnmarshalBinary(data); err != nil {
			t.Error(err)
		}

		dataBefore := make([]byte, len(data))
		copy(dataBefore, data)

		outBefore := make([]byte, len(out))
		copy(outBefore, out)

		out, err := a.AppendBinary(out)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(outBefore, out[:len(outBefore)]) {
			t.Error(outBefore, out)
		}
		if !bytes.Equal(dataBefore, out[len(outBefore):]) {
			t.Error(dataBefore, out[len(outBefore):])
		}
	})
}

func FuzzBitSet256_BinaryEncoding(f *testing.F) {
	f.Add([]byte{0, 1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 0, 0, 0, 0, 0})
	f.Add([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	f.Add([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 32 {
			data = append(data, make([]byte, 32-len(data))...)
		} else if len(data) > 32 {
			data = data[:32]
		}

		var a bitset.BitSet256

		if err := a.UnmarshalBinary(data); err != nil {
			t.Error(err)
		}

		b, err := a.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		if len(b) != 32 {
			t.Error(len(b))
		}

		if !bytes.Equal(data, b) {
			t.Error(data, b)
		}
	})
}

func FuzzBitSet256_GetSet(f *testing.F) {
	f.Add(0, true)
	f.Add(63, true)
	f.Add(64, true)
	f.Add(127, false)

	f.Fuzz(func(t *testing.T, index int, value bool) {
		if index < 0 {
			index = -index
		}

		index = index % bitset.BitSet256Size

		var a bitset.BitSet256

		a.Set(index, value)

		if a.Get(index) != value {
			t.Error(index, a.Get(index), value)
		}

		t.Run("check if other bits remain unchanged unless we're setting/unsetting them", func(t *testing.T) {
			for range 100 {
				if idx := rand.IntN(bitset.BitSet256Size); idx != index {
					v := a.Get(idx)

					a.Set(idx, !v)
					a.Set(idx, v)

					if a.Get(idx) != v {
						t.Error(idx, a.Get(idx), v)
					}
				}
			}
		})
	})
}

func BenchmarkBitSet256_IsZero(b *testing.B) {
	var s bitset.BitSet256

	for b.Loop() {
		s.IsZero()
	}
}

func BenchmarkBitSet256_Set(b *testing.B) {
	var s bitset.BitSet256

	i := rand.IntN(bitset.BitSet256Size)
	v := rand.IntN(2) == 0

	for b.Loop() {
		s.Set(i, v)
	}
}

func BenchmarkBitSet256_Get(b *testing.B) {
	var s bitset.BitSet256

	for i := range bitset.BitSet256Size {
		s.Set(i, rand.IntN(2) == 0)
	}

	i := rand.IntN(bitset.BitSet256Size)

	for b.Loop() {
		s.Get(i)
	}
}

func BenchmarkBitSet256_Bits(b *testing.B) {
	var s bitset.BitSet256

	for i := range bitset.BitSet256Size {
		s.Set(i, rand.IntN(2) == 0)
	}

	for b.Loop() {
		for range bitset.BitSet256Size {
		}
	}
}

func BenchmarkBitSet256_Union(b *testing.B) {
	var x, y bitset.BitSet256

	for i := range bitset.BitSet256Size {
		x.Set(i, rand.IntN(2) == 0)
		y.Set(i, rand.IntN(2) == 0)
	}

	for b.Loop() {
		x.Union(y)
	}
}

func BenchmarkBitSet256_AppendBinary(b *testing.B) {
	var s bitset.BitSet256

	for i := range bitset.BitSet256Size {
		s.Set(i, rand.IntN(2) == 0)
	}

	out := make([]byte, 0, bitset.BitSet256Size*8)

	for b.Loop() {
		s.AppendBinary(out)
	}
}

func BenchmarkBitSet256_MarshalBinary(b *testing.B) {
	var s bitset.BitSet256

	for i := range bitset.BitSet256Size {
		s.Set(i, rand.IntN(2) == 0)
	}

	for b.Loop() {
		s.MarshalBinary()
	}
}

func BenchmarkBitSet256_UnmarshalBinary(b *testing.B) {
	var x, y bitset.BitSet256

	for i := range bitset.BitSet256Size {
		x.Set(i, rand.IntN(2) == 0)
	}

	v, _ := x.MarshalBinary()

	for b.Loop() {
		y.UnmarshalBinary(v)
	}
}
