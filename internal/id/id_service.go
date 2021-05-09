package id

import (
	"math/rand"
	"time"
)

type IdService struct {
}

func New() *IdService {
	return &IdService{}
}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func (i *IdService) randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (i *IdService) NewId(length int) string {

	rand.Seed(time.Now().UnixNano())

	res := i.randSeq(length)
	return res
}
