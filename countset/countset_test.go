package countset_test

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"testing"

	"github.com/ndx-technologies/encoding/countset"
)

func Example_uint64() {
	m := map[uint64]int{
		1234: -11,
	}

	b, _ := countset.MarshalBinary(m)
	m2, _ := countset.UnmarshalBinary[uint64, int](b)

	fmt.Println(maps.Equal(m, m2), m[1234])
	// Output: true -11
}

func Example_customType() {
	type ID uint64
	id := ID(1234)

	m := map[ID]int{
		id: -11,
	}

	b, _ := countset.MarshalBinary(m)
	m2, _ := countset.UnmarshalBinary[ID, int](b)

	fmt.Println(maps.Equal(m, m2), m[id])
	// Output: true -11
}

func FuzzProductCountMap(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Add(32)

	f.Fuzz(func(t *testing.T, n int) {
		n %= 1000

		m := make(map[uint64]int, n)

		for range n {
			m[rand.Uint64()] = rand.Int()
		}

		data, err := countset.MarshalBinary(m)
		if err != nil {
			t.Fatal(err)
		}

		m2, err := countset.UnmarshalBinary[uint64, int](data)
		if err != nil {
			t.Fatal(err)
		}

		if !maps.Equal(m, m2) {
			t.Fatal(m, m2)
		}
	})
}
