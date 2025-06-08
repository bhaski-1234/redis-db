package utils

import (
	"fmt"
	"strconv"
	"time"
)

func EncodeVarIntBigEndian(value int) []byte {
	if value == 0 {
		return []byte{0}
	}

	var groups []byte
	for value > 0 {
		groups = append([]byte{byte(value & 0x7F)}, groups...)
		value >>= 7
	}

	// Set continuation bits for all but the last group
	for i := 0; i < len(groups)-1; i++ {
		groups[i] |= 0x80
	}
	return groups
}

func DecodeVarIntBigEndian(data []byte) uint64 {
	var value uint64
	for i := 0; i < len(data); i++ {
		b := data[i]
		value = (value << 7) | uint64(b&0x7F)
		if b&0x80 == 0 {
			break
		}
	}
	return value
}

func ParseDuration(durationStr string) (time.Duration, error) {
	t, _ := strconv.ParseInt(durationStr, 10, 64)
	if t < 0 {
		return 0, fmt.Errorf("invalid duration: %s", durationStr)
	}
	return time.Duration(t) * time.Second, nil
}
