package base64_test

import (
	"fmt"
	"testing"

	"github.com/ndx-technologies/encoding/base64"
)

func ExampleFloat32Vec() {
	v := base64.Float32Vec{1.0, 2.0, 3.0}
	s, _ := v.MarshalText()

	var u base64.Float32Vec
	_ = u.UnmarshalText(s)

	fmt.Print(string(s), ",", u)
	// Output: AACAPwAAAEAAAEBA,[1 2 3]
}

func ExampleFloat32Vec_empty() {
	var v base64.Float32Vec
	s, _ := v.MarshalText()

	fmt.Print(string(s) == "")
	// Output: true
}

func FuzzFloat32Vec(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		var y base64.Float32Vec
		y.UnmarshalText(b)
	})
}
