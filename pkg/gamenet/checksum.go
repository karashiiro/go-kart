package gamenet

func NetBufferChecksum(buf []byte) uint32 {
	checksum := uint32(0x1234567)
	length := len(buf)
	for i := 0; i < length; i++ {
		checksum += uint32(buf[i]) * uint32(i+1)
	}
	return checksum
}
