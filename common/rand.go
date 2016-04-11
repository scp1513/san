package common

import (
	"math/rand"
	"strconv"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

// 随机字符串
func RandString(size int) string {
	// if size > 1 {
	// 	b := make([]byte, size/2)
	// 	_, err := rand.Read(b)
	// 	if err != nil {
	// 		return ""
	// 	}

	// 	return base64.URLEncoding.EncodeToString(b)
	// } else {
	// 	var b [1]byte
	// 	_, err := rand.Read(b[:])
	// 	if err != nil {
	// 		return ""
	// 	}
	// 	s := base64.URLEncoding.EncodeToString(b[:])
	// 	return string(([]byte(s))[:1])
	// }
	ret := make([]byte, size)

	index := 0
	for {
		randNum := r.Int63()
		s := strconv.FormatInt(randNum, 36)

		copy(ret[index:], s[:])
		index += len(s)
		if index >= size {
			break
		}
	}

	return string(ret[:])
}
