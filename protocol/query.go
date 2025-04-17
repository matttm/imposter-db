package protocol

import (
	"bytes"
	"encoding/binary"
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
	return q
}

func EncodeQuery(q *Query) []byte {
	var b []byte
	w := bytes.NewBuffer(b)
	panic("EncodeQuery not implemented")
	return w.Bytes()
}
func (q *Query) Contains(t string) bool {
	return strings.Contains(q.query, t)
}
