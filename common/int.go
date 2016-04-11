package common

import (
	"fmt"
)

func MakeU64(h, l uint32) uint64 {
	return uint64(h)<<32 | uint64(l)
}

func SplitU64(v uint64) (uint32, uint32) {
	return uint32((v >> 32) & uint64(0x0000FFFF)), uint32(v & uint64(0x0000FFFF))
}

func U32(b []byte) (uint32, error) {
	if len(b) < 4 {
		return 0, fmt.Errorf("不足4个字节 %d", len(b))
	}
	return (uint32(b[3]) << 24) |
		(uint32(b[2]) << 16) |
		(uint32(b[1]) << 8) |
		uint32(b[0]), nil
}

func U32Byte(i uint32) []byte {
	return []byte{
		0: byte(i),
		1: byte(i >> 8),
		2: byte(i >> 16),
		3: byte(i >> 24),
	}
}

func U64(b []byte) (uint64, error) {
	if len(b) < 8 {
		return 0, fmt.Errorf("不足8个字节 %d", len(b))
	}
	return (uint64(b[0])) |
		(uint64(b[1]) << 8) |
		(uint64(b[2]) << 16) |
		(uint64(b[3]) << 24) |
		(uint64(b[4]) << 32) |
		(uint64(b[5]) << 40) |
		(uint64(b[6]) << 48) |
		(uint64(b[7]) << 54), nil
}

func U64Byte(i uint64) []byte {
	return []byte{
		0: byte(i),
		1: byte(i >> 8),
		2: byte(i >> 16),
		3: byte(i >> 24),
		4: byte(i >> 32),
		5: byte(i >> 40),
		6: byte(i >> 48),
		7: byte(i >> 54),
	}
}
