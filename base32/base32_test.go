package base32_test

import (
	"testing"

	"github.com/ndx-technologies/encoding/base32"
)

func TestDecodeError(t *testing.T) {
	s := "-1"
	if _, err := base32.Decode[uint32]([]byte(s)); err == nil {
		t.Error("expected error")
	}
}
