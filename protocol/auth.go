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

// Authentication Methods (Examples - Check your server documentation)
const (
	MYSQL_NATIVE_PASSWORD string = "mysql_native_password"
	CACHING_SHA2_PASSWORD string = "caching_sha2_password"
	MYSQL_CLEAR_PASSWORD  string = "mysql_clear_password"
	SHA256_PASSWORD       string = "sha256_password"
	ED25519               string = "ed25519_plugin"

	// Authentication Stages/Flags (Examples - Check your server documentation)
	AUTH_INITIAL_HANDSHAKE byte = 0x0A // Typically the first packet from the server
	AUTH_MORE_DATA         byte = 0x01 // More data needed for authentication
	AUTH_SWITCH_REQUEST    byte = 0xFE // Server requests authentication method switch
	AUTH_SWITCH_RESPONSE   byte = 0x00 // Client response to authentication switch
)

var authMap map[string]AuthenticationMethod = map[string]AuthenticationMethod{
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_connection_phase_authentication_methods_native_password_authentication.html
	"mysql_native_password": AuthenticationMethod{Fn: sha1Wrapper, Sz: sha1.Size},
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_caching_sha2_authentication_exchanges.html#sect_caching_sha2_info
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
