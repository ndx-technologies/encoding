package base32

// Uint64 is a base32 encoded uint64 in 13 LittleEndian bytes
type Uint64 uint64

func (s Uint64) AppendText(b []byte) ([]byte, error) {
	b = append(
		b,
		digits[(s>>60)&0x0F], // Bits 60–64 (4 bits)
		digits[(s>>55)&0x1F], // Bits 55–59 (5 bits)
		digits[(s>>50)&0x1F], // Bits 50–54 (5 bits)
		digits[(s>>45)&0x1F], // Bits 45–49 (5 bits)
		digits[(s>>40)&0x1F], // Bits 40–44 (5 bits)
		digits[(s>>35)&0x1F], // Bits 35–39 (5 bits)
		digits[(s>>30)&0x1F], // Bits 30–34 (5 bits)
		digits[(s>>25)&0x1F], // Bits 25–29 (5 bits)
		digits[(s>>20)&0x1F], // Bits 20–24 (5 bits)
		digits[(s>>15)&0x1F], // Bits 15–19 (5 bits)
		digits[(s>>10)&0x1F], // Bits 10–14 (5 bits)
		digits[(s>>5)&0x1F],  // Bits 5–9   (5 bits)
		digits[s&0x1F],       // Bits 0–4   (5 bits)
	)
	return b, nil
}

func (s Uint64) MarshalText() ([]byte, error) {
	b := make([]byte, 0, 13)
	return s.AppendText(b)
}

func (s *Uint64) UnmarshalText(b []byte) (err error) {
	if len(b) != 13 {
		return ErrInvalidLength
	}
	v, err := Decode[uint64](b)
	*s = Uint64(v)
	return err
}

func Uint64FromString(s string) (v Uint64, err error) {
	err = v.UnmarshalText([]byte(s))
	return v, err
}

func (s Uint64) String() string {
	b, _ := s.MarshalText()
	return string(b)
}
