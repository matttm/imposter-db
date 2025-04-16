package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
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
func ReadFixedLengthString(r io.Reader, n uint64) string {
	s := make([]byte, n)
	_, err := r.Read(s)
	if err != nil {
		panic(err)
	}
	return string(s)
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

func SaveToFile(data []byte, newDir, newFilename string) error {
	currentTime := time.Now()
	// Create the directory if it doesn't exist
	err := os.MkdirAll(newDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	newFilename = fmt.Sprintf("%s-%s.capture", newFilename, currentTime.Format("2006.01.02 15:04:05"))
	// Construct the full file path
	filePath := filepath.Join(newDir, newFilename)

	// Create the file and write data
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
func xorBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		return nil // Return nil if slices have different lengths
	}
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}
