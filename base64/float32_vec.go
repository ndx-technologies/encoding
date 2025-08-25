package base64

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

var (
	ByteOrder      = binary.LittleEndian
	Base64Encoding = base64.StdEncoding
)

type Float32Vec []float32

func (s Float32Vec) MarshalText() ([]byte, error) {
	if len(s) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	if err := binary.Write(&buf, ByteOrder, []float32(s)); err != nil {
		return nil, err
	}

	return []byte(Base64Encoding.EncodeToString(buf.Bytes())), nil
}

func (s *Float32Vec) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		*s = nil
		return nil
	}

	decoded := make([]byte, Base64Encoding.DecodedLen(len(b)))
	d := base64.NewDecoder(Base64Encoding, bytes.NewReader(b))
	if _, err := d.Read(decoded); err != nil {
		return err
	}

	v := make(Float32Vec, len(decoded)/4)
	r := bytes.NewReader(decoded)
	if err := binary.Read(r, ByteOrder, &v); err != nil {
		return err
	}

	*s = v
	return nil
}
