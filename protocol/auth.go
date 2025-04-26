package protocol

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
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
	FAST_AUTH_SUCCESS      byte = 0x03
	AUTH_SWITCH_RESPONSE   byte = 0x00 // Client response to authentication switch
)

var authMap map[string]AuthenticationMethod = map[string]AuthenticationMethod{
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_connection_phase_authentication_methods_native_password_authentication.html
	MYSQL_NATIVE_PASSWORD: {Fn: sha1Wrapper, Sz: sha1.Size},
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_caching_sha2_authentication_exchanges.html#sect_caching_sha2_info
	CACHING_SHA2_PASSWORD: {Fn: sha256Wrapper, Sz: sha256.Size}, //	this is in development
	SHA256_PASSWORD:       {Fn: sha256Wrapper, Sz: sha256.Size}, //	this is in development
}

func sha1Wrapper(data []byte) []byte {
	sum := sha1.Sum(data)
	return sum[:]
}

func sha256Wrapper(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

// doc: https://dev.mysql.com/blog-archive/preparing-your-community-connector-for-mysql-8-part-2-sha256/
//
//	from doc : It’s important to note that a incompatible change happened in server 8.0.5.  Prior to server 8.0.5 the encryption was done using RSA_PKCS1_PADDING.  With 8.0.5 it is done with RSA_PKCS1_OAEP_PADDING.  This means that if you have implemented support for this authentication scheme for servers prior to 8.0.5 you will need to update your connector to make this change.
func encryptPassword(pemKey, password, salt []byte) []byte {
	log.Printf("PEM KEY: %s", pemKey)
	// Decode PEM to get DER-encoded key
	block, _ := pem.Decode(pemKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		panic("failed to parse PEM public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Print(err)
		panic("RSA parsing error occured")
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("not an RSA public key")
	}
	// XOR (password + '\x00') with the salt (repeated if needed)
	password = append(password, 0x00)
	scrambled := xorScramble(password, salt)
	log.Printf("Scrambled password (hex): %x", scrambled)
	e, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, rsaPub, scrambled, nil)
	log.Printf("Encrypted password len: %d", len(e))
	if err != nil {
		log.Print(err)
		panic("RSA encryption error occured")
	}
	return e
}

// TODO: refactor method to be an enum
func hashPassword(method string, salt []byte, password string) ([]byte, error) {
	if isNonASCIIorEmpty(method) {
		return []byte{}, fmt.Errorf("Authentication method is undecipherable")
	}
	if len(salt) > 20 {
		salt = salt[:20]
	}
	log.Printf("Hashing password with %x", salt)
	log.Printf("Hashing password %s via %s method", password, method)
	if authMeth, ok := authMap[method]; ok {
		var scrambled []byte
		// https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_connection_phase_authentication_methods_native_password_authentication.html
		if method == "mysql_native_password" {
			// https://dev.mysql.com/doc/dev/mysql-server/8.0.40/page_protocol_connection_phase_authentication_methods_native_password_authentication.html
			stage1 := authMeth.Fn([]byte(password))
			dub := authMeth.Fn(stage1[:])
			stage2 := authMeth.Fn(append(salt, dub[:]...))

			scrambled = make([]byte, authMeth.Sz)
			for i := 0; i < authMeth.Sz; i++ {
				scrambled[i] = stage1[i] ^ stage2[i]
			}
		}
		// https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_caching_sha2_authentication_exchanges.html#sect_caching_sha2_info
		if method == "caching_sha2_password" || method == SHA256_PASSWORD {
			stage1 := authMeth.Fn([]byte(password))
			dub := authMeth.Fn(stage1[:])
			stage2 := authMeth.Fn(append(dub[:], salt...))

			scrambled = make([]byte, authMeth.Sz)
			for i := 0; i < authMeth.Sz; i++ {
				scrambled[i] = stage1[i] ^ stage2[i]
			}
		}
		return scrambled, nil
	}
	return []byte{}, fmt.Errorf("Unknown authentication method: %s", method)
}
