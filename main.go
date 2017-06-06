package main

import "fmt"
import (
	"time"
	"math/rand"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var runesForID = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetUUID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runesForID[rand.Intn(len(runesForID))]
	}
	return string(b)
}

func main()  {

}