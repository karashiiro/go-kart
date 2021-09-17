package gamenet

import (
	"math"
	"unsafe"

	"github.com/karashiiro/gokart/pkg/network"
	"github.com/tav/golly/lzf"
)

// Whereas the original implementation uses a linked list of file
// transactions, we just send the files one by one in a blocking
// manner, since we have the benefit of goroutines.

func SendFileMemoryCompressed(conn network.Connection, data []byte) error {
	compressed := lzf.Compress(data)
	return SendFileMemory(conn, compressed)
}

func SendFileMemory(conn network.Connection, data []byte) error {
	pos := 0
	size := len(data)

	for pos < size {
		nextChunkSize := int(math.Min(float64(size), 1011))

		f := &FileTxPak{
			FileId:   uint8(uintptr(unsafe.Pointer(&data)) >> 8),
			Position: uint32(pos),
			Size:     uint16(nextChunkSize),
		}

		for i := 0; i < nextChunkSize; i++ {
			f.Data[i] = data[pos]
			pos++
		}

		err := SendPacket(conn, f)
		if err != nil {
			return err
		}
	}

	return nil
}
