package bitset

import (
	"encoding/binary"
	"errors"
)

const BitSet256Size = 256

type BitSet256 struct{ words [4]uint64 }

func (s BitSet256) IsZero() bool { return s == BitSet256{} }

func (s *BitSet256) Set(i int, v bool) {
	if i < 0 || i >= BitSet256Size {
		return
	}
	word, bit := i/64, uint64(i%64)
	if v {
		s.words[word] |= (1 << bit)
	} else {
		s.words[word] &^= (1 << bit)
	}
}

func (s BitSet256) Get(i int) bool {
	if i < 0 || i >= BitSet256Size {
		return false
	}
	word, bit := i/64, uint64(i%64)
	return (s.words[word] & (1 << bit)) != 0
}

func (s *BitSet256) Union(other BitSet256) {
	s.words[0] |= other.words[0]
	s.words[1] |= other.words[1]
	s.words[2] |= other.words[2]
	s.words[3] |= other.words[3]
}

func (s BitSet256) AppendBinary(b []byte) ([]byte, error) {
	for _, w := range s.words {
		b = binary.LittleEndian.AppendUint64(b, w)
	}
	return b, nil
}

func (s BitSet256) MarshalBinary() ([]byte, error) { return s.AppendBinary(make([]byte, 0, 32)) }

func (s *BitSet256) UnmarshalBinary(data []byte) error {
	if len(data) != 32 {
		return errors.New("invalid length")
	}
	for i := range 4 {
		s.words[i] = binary.LittleEndian.Uint64(data[i*8 : (i+1)*8])
	}
	return nil
}
