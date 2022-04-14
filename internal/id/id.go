package id

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type (
	IdServiceProvider interface {
		IdService() IdService
	}

	IdService interface {
		NewID(length int) string
		NewUUID() (uuid.UUID, error)
	}

	idService struct {
	}
)

func NewIdService() *idService {
	return &idService{}
}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func (i *idService) randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (i *idService) NewID(length int) string {

	rand.Seed(time.Now().UnixNano())

	res := i.randSeq(length)
	return res
}

func (i *idService) NewUUID() (uuid.UUID, error) {
	return uuid.NewRandom()
}
