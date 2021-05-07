package services

import (
	"math/rand"
	"time"
)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewId() string {

	rand.Seed(time.Now().UnixNano())

	res := randSeq(16)
	return res
}
