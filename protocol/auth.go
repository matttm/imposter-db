package protocol

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
)

type AuthenticationMethod struct {
	Fn func([]byte) []byte
	Sz int
}

var authMap map[string]AuthenticationMethod = map[string]AuthenticationMethod{
	"mysql_native_password": AuthenticationMethod{Fn: sha1Wrapper, Sz: sha1.Size},
	"caching_sha2_password": AuthenticationMethod{Fn: sha256Wrapper, Sz: sha256.Size}, //	this is in development
}

func sha1Wrapper(data []byte) []byte {
	sum := sha1.Sum(data)
	return sum[:]
}

func sha256Wrapper(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

func encryptPassword(pubKey, password []byte) []byte {
	pub, err := x509.ParsePKCS1PublicKey(pubKey)
	if err != nil {
		panic("RSA encryption error occured")
	}
	e, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, password, nil)
	if err != nil {
		panic("RSA encryption error occured")
	}
	return e
}
