package protocol

import (
	"bytes"
	"encoding/binary"
	"log"
	"regexp"
	"strings"
)

type Query struct {
	command       uint8
	paramCount    uint64
	paramSetCount uint64 // always 1 rn
	// other fields...
	query string
}

func DecodeQuery(flags uint32, b []byte) *Query {
	q := &Query{}
	r := bytes.NewReader(b)
	err := binary.Read(r, binary.LittleEndian, &q.command)
	if err != nil {
		panic(err)
	}
	if flags&CLIENT_QUERY_ATTRIBUTES != 0 {
		q.paramCount, _ = ReadVarLengthInt(r)
		q.paramSetCount, _ = ReadVarLengthInt(r)
	}
	q.query = ReadFixedLengthString(r, uint64(r.Len()))
	q.query = strings.ToLower(q.query)
	return q
}

func EncodeQuery(q *Query) []byte {
	var b []byte
	w := bytes.NewBuffer(b)
	panic("EncodeQuery not implemented")
	return w.Bytes()
}
func (q *Query) Contains(t string) bool {
	log.Printf("Checking for '%s' in '%s'", t, q.query)
	b, err := regexp.MatchString(`from `+regexp.QuoteMeta(t)+`\s+`, q.query)
	if err != nil {
		panic(err)
	}
	return b
}
