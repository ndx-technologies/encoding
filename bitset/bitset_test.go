package bitset_test

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/ndx-technologies/encoding/bitset"
)

func ExampleBitSet_empty() {
	var empty bitset.BitSet
	fmt.Print(empty.IsZero())
	// Output: true
}

func ExampleBitSet_Get() {
	var a bitset.BitSet
	a.Set(11, true)
	a.Set(1000, true)

	fmt.Print(a.IsZero(), a.Get(11), a.Get(64), a.Get(256))
	// Output: false true false false
}

func ExampleBitSet_Union() {
	var a, b bitset.BitSet

	a.Set(1, true)
	b.Set(1000, true)

	a.Union(b)

	fmt.Print(a.Get(0), a.Get(1), a.Get(2), a.Get(1000), a.Size())
	// Output: false true false true 1024
}

func TestBitSet_encodingError(t *testing.T) {
	var a bitset.BitSet
	if err := a.UnmarshalBinary([]byte{0, 1, 2, 3, 4, 5}); err == nil {
		t.Error("expected error")
	}
}

func FuzzBitSet_AppendBinary(f *testing.F) {
	f.Add([]byte(nil), []byte{1, 2, 3, 4, 5, 6, 7, 8})
	f.Add([]byte{1, 2, 3}, []byte{1, 2, 3, 4, 5, 6, 7, 8})

	f.Fuzz(func(t *testing.T, out, data []byte) {
		if len(data) < 16 {
			data = append(data, make([]byte, 16-len(data))...)
		} else if len(data) > 16 {
			data = data[:16]
		}

		var a bitset.BitSet

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

func FuzzBitSet_BinaryEncoding(f *testing.F) {
	f.Add([]byte{0, 1, 2, 3, 4, 5, 6, 7})
	f.Add([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	f.Add([]byte{1, 0, 0, 0, 0, 0, 0, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 16 {
			data = append(data, make([]byte, 16-len(data))...)
		} else if len(data) > 16 {
			data = data[:16]
		}

		var bs bitset.BitSet

		if err := bs.UnmarshalBinary(data); err != nil {
			t.Error(err)
		}

		b, err := bs.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(data, b) {
			t.Error(data, b)
		}
	})
}

func FuzzBitSet_GetSet(f *testing.F) {
	f.Add(0, true)
	f.Add(63, true)
	f.Add(64, true)
	f.Add(127, false)

	f.Fuzz(func(t *testing.T, index int, value bool) {
		var bs bitset.BitSet

		if index < 0 {
			index = -index
		}

		index = index % (1024 * 64)

		bs.Set(index, value)

		if bs.Get(index) != value {
			t.Error(index, bs.Get(index), value)
		}

		t.Run("check if other bits remain unchanged unless we're setting/unsetting them", func(t *testing.T) {
			for range 100 {
				if idx := rand.IntN(128); idx != index {
					v := bs.Get(idx)

					bs.Set(idx, !v)
					bs.Set(idx, v)

					if bs.Get(idx) != v {
						t.Error(idx, bs.Get(idx), v)
					}
				}
			}
		})
	})
}

func BenchmarkBitSet_IsEmpty(b *testing.B) {
	var s bitset.BitSet

	for b.Loop() {
		s.IsZero()
	}
}

func BenchmarkBitSet_Set(b *testing.B) {
	for _, n := range []int{256, 700} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			var s bitset.BitSet

			i := rand.IntN(n)
			v := rand.IntN(2) == 0

			for b.Loop() {
				s.Set(i, v)
			}
		})
	}
}

func BenchmarkBitSet_Get(b *testing.B) {
	for _, n := range []int{256, 700} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			var s bitset.BitSet

			for i := range n {
				s.Set(i, rand.IntN(2) == 0)
			}

			i := rand.IntN(n)

			for b.Loop() {
				s.Get(i)
			}
		})
	}
}

func BenchmarkBitSet_Bits(b *testing.B) {
	for _, n := range []int{256, 700} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			var s bitset.BitSet

			for i := range n {
				s.Set(i, rand.IntN(2) == 0)
			}

			for b.Loop() {
				for i := 0; i < s.Size(); i++ {
				}
			}
		})
	}
}

func BenchmarkBitSet_Union(b *testing.B) {
	for _, n := range []int{256, 700} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			var x, y bitset.BitSet

			for i := range n {
				x.Set(i, rand.IntN(2) == 0)
				y.Set(i, rand.IntN(2) == 0)
			}

			for b.Loop() {
				x.Union(y)
			}
		})
	}
}

func BenchmarkBitSet_AppendBinary(b *testing.B) {
	for _, n := range []int{256, 700} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			var s bitset.BitSet

			for i := range n {
				s.Set(i, rand.IntN(2) == 0)
			}

			out := make([]byte, 0, n*2*8)

			for b.Loop() {
				s.AppendBinary(out)
			}
		})
	}
}

func BenchmarkBitSet_MarshalBinary(b *testing.B) {
	for _, n := range []int{256, 700} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			var s bitset.BitSet

			for i := range n {
				s.Set(i, rand.IntN(2) == 0)
			}

			for b.Loop() {
				s.MarshalBinary()
			}
		})
	}
}

func BenchmarkBitSet_UnmarshalBinary(b *testing.B) {
	for _, n := range []int{256, 700} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			var x, y bitset.BitSet

			for i := range n {
				x.Set(i, rand.IntN(2) == 0)
			}

			v, _ := x.MarshalBinary()

			for b.Loop() {
				y.UnmarshalBinary(v)
			}
		})
	}
}
