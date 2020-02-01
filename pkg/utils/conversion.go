package utils

import "encoding/binary"

// Int64ToBytes returns an 8-byte big endian representation of v.
func Int64ToBytes(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// BytesToInt64 return an int64 of v.
func BytesToInt64(v []byte) int64 {
	return int64(binary.BigEndian.Uint64(v))
}
