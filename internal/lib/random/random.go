package random

import (
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewRandomString(length int) string {
	b:= make([]byte, length)
	for i:=range b{
		b[i]= letterBytes[random.Intn(len(letterBytes))]
	}
	return string(b)
}