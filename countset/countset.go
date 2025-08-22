package countset

import (
	"encoding/binary"
	"errors"
)

var ByteOrder = binary.LittleEndian

func AppendBinary[K ~uint64, V ~int](b []byte, mp map[K]V) ([]byte, error) {
	b = ByteOrder.AppendUint64(b, uint64(len(mp)))
	for k, v := range mp {
		b = ByteOrder.AppendUint64(b, uint64(k))
		b = ByteOrder.AppendUint64(b, uint64(v))
	}
	return b, nil
}

func MarshalBinary[K ~uint64, V ~int](mp map[K]V) ([]byte, error) {
	return AppendBinary(make([]byte, 0, 8+len(mp)*16), mp)
}

func UnmarshalBinary[K ~uint64, V ~int](data []byte) (map[K]V, error) {
	if len(data) < 1 {
		return nil, errors.New("invalid length")
	}

	count := int(ByteOrder.Uint64(data[:8]))
	if count == 0 {
		return nil, nil
	}
	data = data[8:]

	if len(data) != count*16 {
		return nil, errors.New("invalid length")
	}

	mp := make(map[K]V, count)

	for range count {
		k, v := ByteOrder.Uint64(data[:8]), ByteOrder.Uint64(data[8:16])
		data = data[16:]
		mp[K(k)] = V(v)
	}

	return mp, nil
}
