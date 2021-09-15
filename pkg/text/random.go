package text

import "math/rand"

const LETTERBYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTERBYTES[rand.Intn(len(LETTERBYTES))]
	}
	return string(b)
}
