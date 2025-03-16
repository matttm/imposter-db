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
func ParseVarLengthInt(r io.Reader) (uint64, int) {
	var x byte
	var ret uint64
	var sz int
	err := binary.Read(r, binary.LittleEndian, x)
	if err != nil {
		panic(err)
	}
	switch x {
	case 0xFE:
		var a uint64
		err = binary.Read(r, binary.LittleEndian, a)
		ret = a
		sz = 8
		break
	case 0xFD:
		panic("parsing 3-byte var length int not implemented")
	case 0xFC:
		var c uint16
		err = binary.Read(r, binary.LittleEndian, c)
		ret = uint64(c)
		sz = 2
		break
	default:
		return uint64(x), 1
	}
	if err != nil {
		panic(err)
	}
	return ret, sz + 1 // we add 1 to account for the examined byte
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
		if err != nil || b == 0x0 {
			break
		}
		s = append(s, b)
		n--
	}
	return string(s)
}
