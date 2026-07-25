package state

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/matttm/imposter-db/protocol"
)

// AuthMoreDataState represents the multi-step authentication phase
// The server sends authentication data and waits for client responses
// This state is used when authentication requires multiple exchanges
type AuthMoreDataState struct {
	nonce []byte
}

func (amd *AuthMoreDataState) IsComplete() bool {
	// AuthMoreDataState is not an exit state
	return false
}

func (amd *AuthMoreDataState) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) (State, error) {
	log.Println("AuthMoreDataState: Processing auth more data...")

	// Read the AUTH_MORE_DATA packet from server
	b, err := protocol.ReadPacket(remote)
	if err != nil {
		return nil, fmt.Errorf("AuthMoreDataState: error reading AUTH_MORE_DATA: %w", err)
	}

	if len(b) < 5 || b[4] != protocol.AUTH_MORE_DATA {
		return nil, fmt.Errorf("AuthMoreDataState: expected AUTH_MORE_DATA packet but got %x", b[4])
	}

	// Check the auth more data type (at index 5)
	if len(b) > 5 {
		authType := b[5]

		// FAST_AUTH_SUCCESS case
		if authType == protocol.FAST_AUTH_SUCCESS {
			log.Println("AuthMoreDataState: FAST_AUTH_SUCCESS received")
			if client != nil {
				_, err := client.Write(b)
				if err != nil {
					return nil, fmt.Errorf("AuthMoreDataState: error sending to client: %w", err)
				}
			}

			// Read next packet (should be OK packet)
			b, err = protocol.ReadPacket(remote)
			if err != nil {
				return nil, fmt.Errorf("AuthMoreDataState: error reading after FAST_AUTH_SUCCESS: %w", err)
			}

			if !protocol.IsOkPacket(b) {
				return nil, fmt.Errorf("AuthMoreDataState: expected OK packet after FAST_AUTH_SUCCESS but got %x", b[4])
			}

			log.Println("AuthMoreDataState: OK packet received after FAST_AUTH_SUCCESS")
			if client != nil {
				_, err := client.Write(b)
				if err != nil {
					return nil, fmt.Errorf("AuthMoreDataState: error sending OK packet to client: %w", err)
				}
			}
			return &ReadyState{}, nil
		}

		// PERFORM_FULL_AUTH case
		if authType == protocol.PERFORM_FULL_AUTH {
			log.Println("AuthMoreDataState: PERFORM_FULL_AUTH requested")

			// Request public key from server
			log.Println("AuthMoreDataState: Requesting server's public key")
			reqKeyPacket := protocol.PackPayload([]byte{0x02}, b[3]+1)
			_, err := remote.Write(reqKeyPacket)
			if err != nil {
				return nil, fmt.Errorf("AuthMoreDataState: error requesting public key: %w", err)
			}

			// Read public key packet
			pemPacket, err := protocol.ReadPacket(remote)
			if err != nil {
				return nil, fmt.Errorf("AuthMoreDataState: error reading public key packet: %w", err)
			}

			// Extract PEM from packet (skip header and auth more data marker)
			pem := pemPacket[4:] // remove packet header
			if len(pem) > 0 && pem[0] == protocol.AUTH_MORE_DATA {
				pem = pem[1:] // remove auth more data marker
			}

			// Encrypt password with public key
			encryptedPwd := protocol.EncryptPassword(pem, []byte(password), amd.nonce)
			encryptedPacket := protocol.PackPayload(encryptedPwd, pemPacket[3]+1)

			// Send encrypted password to server
			_, err = remote.Write(encryptedPacket)
			if err != nil {
				return nil, fmt.Errorf("AuthMoreDataState: error sending encrypted password: %w", err)
			}

			// Read final response
			b, err = protocol.ReadPacket(remote)
			if err != nil {
				return nil, fmt.Errorf("AuthMoreDataState: error reading final auth response: %w", err)
			}

			if protocol.IsOkPacket(b) {
				log.Println("AuthMoreDataState: OK packet received after full auth")
				if client != nil {
					_, err := client.Write(b)
					if err != nil {
						return nil, fmt.Errorf("AuthMoreDataState: error sending OK packet to client: %w", err)
					}
				}
				return &ReadyState{}, nil
			}

			return nil, fmt.Errorf("AuthMoreDataState: expected OK packet but got %x", b[4])
		}
	}

	return nil, fmt.Errorf("AuthMoreDataState: unexpected auth more data type")
}
