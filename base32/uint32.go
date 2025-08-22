package base32

// Uint32 is a base32 encoded uint32 in 7 LittleEndian bytes
type Uint32 uint32

func (s Uint32) AppendText(b []byte) ([]byte, error) {
	b = append(
		b,
		digits[(s>>30)&0x03], // Bits 30-32 (2 bits)
		digits[(s>>25)&0x1F], // Bits 25–29 (5 bits)
		digits[(s>>20)&0x1F], // Bits 20–24 (5 bits)
		digits[(s>>15)&0x1F], // Bits 15–19 (5 bits)
		digits[(s>>10)&0x1F], // Bits 10–14 (5 bits)
		digits[(s>>5)&0x1F],  // Bits 5–9   (5 bits)
		digits[s&0x1F],       // Bits 0–4   (5 bits)
	)
	return b, nil
}

func (s Uint32) MarshalText() ([]byte, error) { return s.AppendText(make([]byte, 0, 7)) }

func (s *Uint32) UnmarshalText(b []byte) (err error) {
	if len(b) != 7 {
		return ErrInvalidLength
	}
	v, err := Decode[uint32](b)
	*s = Uint32(v)
	return err
}

func Uint32FromString(s string) (v Uint32, err error) {
	err = v.UnmarshalText([]byte(s))
	return v, err
}

func (s Uint32) String() string {
	b, _ := s.MarshalText()
	return string(b)
}
