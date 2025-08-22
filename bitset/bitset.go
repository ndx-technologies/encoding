package bitset

import (
	"encoding/binary"
	"errors"
)

type BitSet struct{ words []uint64 }

func (s BitSet) Size() int { return len(s.words) * 64 }

func (s BitSet) IsZero() bool {
	for _, w := range s.words {
		if w != 0 {
			return false
		}
	}
	return true
}

func (s *BitSet) Set(i int, v bool) {
	word, bit := i/64, uint64(i%64)
	if word >= len(s.words) {
		if v {
			s.words = append(s.words, make([]uint64, word-len(s.words)+1)...)
		} else {
			return
		}
	}
	if v {
		s.words[word] |= 1 << bit
	} else {
		s.words[word] &^= 1 << bit
	}
}

func (s BitSet) Get(i int) bool {
	word, bit := i/64, uint64(i%64)
	return word < len(s.words) && ((s.words[word] & (1 << bit)) != 0)
}

func (s *BitSet) Union(other BitSet) {
	if len(s.words) < len(other.words) {
		s.words = append(s.words, make([]uint64, len(other.words)-len(s.words))...)
	}
	for i, w := range other.words {
		s.words[i] |= w
	}
}

func (s BitSet) AppendBinary(b []byte) ([]byte, error) {
	for _, w := range s.words {
		b = binary.LittleEndian.AppendUint64(b, w)
	}
	return b, nil
}

func (s BitSet) MarshalBinary() ([]byte, error) {
	return s.AppendBinary(make([]byte, 0, 8*len(s.words)))
}

func (s *BitSet) UnmarshalBinary(data []byte) error {
	if len(data)%8 != 0 {
		return errors.New("invalid length")
	}
	if n := len(data) / 8; cap(s.words) < n {
		s.words = make([]uint64, n)
	}
	for i := 0; i < len(s.words); i++ {
		s.words[i] = binary.LittleEndian.Uint64(data[i*8 : (i+1)*8])
	}
	return nil
}
