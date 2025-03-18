package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
)

func Max(a, b uint8) uint8 {
	if a > b {
		return a
	} else {
		return b
	}
}

// Function ParseVarLengthInt will parse bytes into an int and
//
// Doc fof vaf length int: https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_basic_dt_integers.html#sect_protocol_basic_dt_int_le
//
// return the integer and the number of bytes the int is, in mem
func ReadVarLengthInt(r io.Reader) (uint64, int) {
	var x byte
	var sz int
	b := make([]byte, 1)
	_, err := r.Read(b)
	x = b[0]
	if err != nil {
		panic(err)
	}
	switch x {
	case 0xFE:
		sz = 8
		break
	case 0xFD:
		sz = 3
	case 0xFC:
		sz = 2
		break
	default:
		return uint64(x), 1
	}
	if err != nil {
		panic(err)
	}
	buf := make([]byte, sz)
	return binary.LittleEndian.Uint64(buf), sz + 1 // we add 1 to account for the examined byte
}
func ReadNBytes(r bytes.Reader, n uint) []byte {
	b := make([]byte, n)
	_, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// Function ReadNullTerminatedString
//
// returns a string, and increments reader's position to after string's null byte
func ReadNullTerminatedString(r *bytes.Reader) string {
	name := []byte{}
	for true {
		b, err := r.ReadByte()
		if err != nil || b == 0x0 {
			break
		}
		name = append(name, b)
	}
	return string(name)
}
func ReadFixedLengthString(r *bytes.Reader, n uint) string {
	s := []byte{}
	for n > 0 {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		s = append(s, b)
		n--
	}
	return string(s)
}
