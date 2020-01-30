package utils

import "encoding/binary"

// I64tob returns an 8-byte big endian representation of v.
func I64tob(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Btoi64 return an int64 of v.
func Btoi64(v []byte) int64 {
	return int64(binary.BigEndian.Uint64(v))
}
