package base10

import (
	"errors"
)

var ErrL9InvalidFormat = errors.New("l9: invalid format")

const (
	MaxL9 = 999_999_999
	MinL9 = 0
)

// L9 is base10 chars encoding least significant 9 decimals of uint32.
type L9 struct{ v uint32 }

func L9FromUint32(v uint32) L9 { return L9{v: v} }

func L9FromString(s string) (L9, error) {
	var v L9
	err := (&v).UnmarshalText([]byte(s))
	return v, err
}

func (s L9) UInt32() uint32 { return s.v }

func (s L9) IsZero() bool { return s.v == 0 }

func (s L9) AppendText(b []byte) ([]byte, error) {
	n := len(b)
	b = append(b, make([]byte, 9)...)

	v := s.v
	v, b[n+8] = v/10, '0'+byte(v%10)
	v, b[n+7] = v/10, '0'+byte(v%10)
	v, b[n+6] = v/10, '0'+byte(v%10)
	v, b[n+5] = v/10, '0'+byte(v%10)
	v, b[n+4] = v/10, '0'+byte(v%10)
	v, b[n+3] = v/10, '0'+byte(v%10)
	v, b[n+2] = v/10, '0'+byte(v%10)
	v, b[n+1] = v/10, '0'+byte(v%10)
	_, b[n+0] = v/10, '0'+byte(v%10)
	return b, nil
}

func (s L9) MarshalText() ([]byte, error) { return s.AppendText(make([]byte, 0, 9)) }

func (s *L9) UnmarshalText(b []byte) error {
	if len(b) != 9 {
		return ErrL9InvalidFormat
	}
	s.v = 0
	for i := 0; i < 9; i++ {
		if b[i] < '0' || b[i] > '9' {
			return ErrL9InvalidFormat
		}
		s.v *= 10
		s.v += uint32(b[i] - '0')
	}
	return nil
}

func (s L9) String() string {
	b, _ := s.MarshalText()
	return string(b)
}
