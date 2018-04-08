package p2p

import (
	"fmt"
	"io"
)

func checkEqual(ref, current []byte) bool {
	for i := 0; i < len(current); i++ {
		if ref[i] != current[i] {
			return false
		}
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func uint32ToBytes(buf []byte, n uint32) {
	buf[0] = byte(n >> 24)
	buf[1] = byte(n >> 16)
	buf[2] = byte(n >> 8)
	buf[3] = byte(n)
}

func bytesToUint32(buf []byte) uint32 {
	return (uint32(buf[0]) << 24) |
		(uint32(buf[1]) << 16) |
		(uint32(buf[2]) << 8) | uint32(buf[3])
}

func writeNBOUint32(w io.Writer, n uint32) (err error) {
	buf := make([]byte, 4)
	uint32ToBytes(buf, n)
	_, err = w.Write(buf[0:])
	return
}

func readNBOUint32(r io.Reader) (n uint32, err error) {
	var buf [4]byte
	_, err = r.Read(buf[0:])
	if err != nil {
		return
	}
	n = bytesToUint32(buf[0:])
	return
}

func humanSize(value float64) string {
	switch {
	case value > 1<<30:
		return fmt.Sprintf("%.2f GB", value/(1<<30))
	case value > 1<<20:
		return fmt.Sprintf("%.2f MB", value/(1<<20))
	case value > 1<<10:
		return fmt.Sprintf("%.2f kB", value/(1<<10))
	}
	return fmt.Sprintf("%.2f B", value)
}
