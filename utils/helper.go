package utils

import (
	"math/rand"
	"time"
)

const codeLength = 8
const codeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCode(codes []string) (string, error) {

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	existingCodes := make(map[string]bool)
	for _, code := range codes {
		existingCodes[code] = true
	}

	for {
		b := make([]byte, codeLength)
		for i := range b {
			b[i] = codeChars[r.Intn(len(codeChars))]
		}
		newCode := string(b)
		if !existingCodes[newCode] {
			return newCode, nil
		}
	}
}
