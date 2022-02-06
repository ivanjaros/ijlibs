package random

import "crypto/rand"

func CryptoBytes(ln int) []byte {
	secret := make([]byte, ln)
	_, err := rand.Read(secret)
	if err != nil {
		panic(err)
	}
	return secret
}

func CryptoString(ln int) string {
	return string(CryptoBytes(ln))
}
