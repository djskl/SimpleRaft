package utils

import (
	"time"
	"strconv"
	"math/rand"
)

var runesForID = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomStr(n int) string {
	if n < 1 {
		return ""
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = runesForID[rand.Intn(len(runesForID))]
	}
	return string(b)
}

func GetUUID(n int) string {
	if n < 1 {
		return ""
	}

	time_str := strconv.FormatInt(time.Now().UnixNano(), 16)
	size := len(time_str)
	if n < size {
		x1, x2 := n/2, n-n/2
		return time_str[size-x1:] + randomStr(x2)
	}

	return time_str + randomStr(n-size)

}
