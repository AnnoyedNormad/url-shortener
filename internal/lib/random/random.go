package random

import (
	"math/rand"
	"time"
)

const alphabet = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func NewRandomString(charsQty int) string {
	var genString string
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < charsQty; i++ {
		randCharI := rnd.Intn(len(alphabet) - 1)
		randChar := alphabet[randCharI]

		genString += string(randChar)
	}

	return genString
}
