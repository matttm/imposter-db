package protocol

import (
	"crypto/sha1"
	"crypto/sha256"
)

type AuthenticationMethod struct {
	Fn func([]byte) []byte
	Sz int
}

var authMap map[string]AuthenticationMethod = map[string]AuthenticationMethod{
	"mysql_native_password": AuthenticationMethod{Fn: sha1Wrapper, Sz: sha1.Size},
	"caching_sha2_password": AuthenticationMethod{Fn: sha256Wrapper, Sz: sha256.Size},
}

func sha1Wrapper(data []byte) []byte {
	sum := sha1.Sum(data)
	return sum[:]
}

func sha256Wrapper(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}
