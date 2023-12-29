package random

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

var (
	ErrAliasesAreOver = errors.New("aliases are over")
)

func NewRandomString(size int, alphabet string) string {
	if len(alphabet) == 1 {
		return strings.Repeat(alphabet, size)
	}
	var genString string
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < size; i++ {
		randCharI := rnd.Intn(len(alphabet))
		randChar := alphabet[randCharI]

		genString += string(randChar)
	}

	return genString
}
