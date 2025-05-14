package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"
)

var NULL = []byte{0x0}

func Max(a, b uint8) uint8 {
	if a > b {
		return a
	} else {
		return b
	}
}

// Function ParseVarLengthInt will parse bytes into an int and
//
// Doc for var length int: https://dev.mysql.com/doc/dev/mysql-server/8.4.3/page_protocol_basic_dt_integers.html#sect_protocol_basic_dt_int_le
//
// return the integer and the number of bytes the int is, in mem
func ReadVarLengthInt(r io.Reader) (uint64, int) {
	var x byte = ReadByte(r)
	var sz int
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
	buf := make([]byte, sz)
	_, err := r.Read(buf)
	if err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint64(buf), sz + 1 // we add 1 to account for the examined byte
}

// ReadLengthEncodedString reads a length-encoded string from the provided io.Reader.
// It first reads a variable-length integer to determine the length of the string,
// then reads and returns the string of that length.
// If an error occurs during reading, an empty string is returned.
func ReadLengthEncodedString(r io.Reader) string {
	n, _ := ReadVarLengthInt(r)
	return ReadFixedLengthString(r, n)
}

// Read N bytes while preserving edndian-ness
func ReadNBytes(r io.Reader, n uint) []byte {
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
func ReadNullTerminatedString(r io.Reader) string {
	name := []byte{}
	for true {
		b := ReadByte(r)
		if b == 0x0 {
			break
		}
		name = append(name, b)
	}
	return string(name)
}

// ReadFixedLengthString reads exactly n bytes from the provided io.Reader and returns them as a string.
// If an error occurs during reading, the function panics.
// Note: The returned string may contain null bytes or other non-printable characters if present in the input.
func ReadFixedLengthString(r io.Reader, n uint64) string {
	s := make([]byte, n)
	_, err := r.Read(s)
	if err != nil {
		panic(err)
	}
	return string(s)
}

// ReadStringEOF reads and returns a string from the provided bytes.Reader until EOF.
// It utilizes ReadFixedLengthString to read the remaining bytes as a string.
func ReadStringEOF(r *bytes.Reader) string {
	return ReadFixedLengthString(r, uint64(r.Len()))
}

func ReadByte(r io.Reader) byte {
	b := make([]byte, 1)
	_, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	return b[0]
}

// Function ReadNullTerminatedString
//
// returns a string, and increments reader's position to after string's null byte
func WriteNullTerminatedString(w io.Writer, s string) {
	err := binary.Write(w, binary.LittleEndian, []byte(s))
	if err != nil {
		panic(err)
	}
	err = binary.Write(w, binary.LittleEndian, NULL)
	if err != nil {
		panic(err)
	}
}
func WriteLengthEncodedString(w io.Writer, s string) {
	err := binary.Write(w, binary.LittleEndian, uint8(len(s)))
	if err != nil {
		panic(err)
	}
	err = binary.Write(w, binary.LittleEndian, []byte(s))
	if err != nil {
		panic(err)
	}
}

func customPanic(err error) {
	if err == io.EOF {
		fmt.Println("Client disconnected")
		return
	}
	panic(err)
}

func isNonASCIIorEmpty(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, r := range s {
		if r > utf8.RuneSelf {
			return true
		}
	}
	return false
}

// SaveToFile writes the provided data to a file in the specified directory with a timestamped filename.
// It creates the directory if it does not exist. The filename is constructed by appending the current
// timestamp to the provided newFilename, followed by the ".capture" extension. Returns an error if
// directory creation, file creation, or writing fails.
//
// Parameters:
//   - data: The byte slice to be written to the file.
//   - newDir: The directory where the file will be saved.
//   - newFilename: The base name for the file (timestamp and extension will be appended).
//
// Returns:
//   - error: An error if any operation fails, otherwise nil.
