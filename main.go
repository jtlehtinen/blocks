package main

import (
	"crypto/sha256"
)

func DoubleSha256(b []byte) []byte {
	sha := sha256.New()
	sha.Write(b)
	intermediate := sha.Sum(nil)

	sha.Reset()
	sha.Write(intermediate)
	return sha.Sum(nil)
}

func main() {

}
