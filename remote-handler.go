package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
)

type ProxyRequestHandler struct {
	remoteDb *client.Conn
	spoof    *client.Conn
	spoofed  string
}

func NewRemoteHandler(c *client.Conn, t string, db *client.Conn) ProxyRequestHandler {
	return ProxyRequestHandler{remoteDb: c, spoofed: t, spoof: db}
}

// UseDB is called for COM_INIT_DB
func (h ProxyRequestHandler) UseDB(dbName string) error {
	log.Println("In UseDb")
	return h.remoteDb.UseDB(dbName)
}

// HandleQuery is called for COM_QUERY
func (h ProxyRequestHandler) HandleQuery(query string) (*mysql.Result, error) {
	log.Println("In HandleQuery")
	log.Println(query)
	// SET doesn't seem to work due to EOF/OK differences
	if strings.Contains(query, `SET NAMES`) {
		return nil, nil
	}
	if strings.Contains(query, `set autocommit=0`) {
		return nil, nil
	}
	if strings.Contains(query, h.spoofed) {
		r, err := h.spoof.Execute(query)
		return r, err
	}
	res, err := h.remoteDb.Execute(query)

	return res, err
}

// HandleFieldList is called for COM_FIELD_LIST packets
// Note that COM_FIELD_LIST has been deprecated since MySQL 5.7.11
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_field_list.html
func (h ProxyRequestHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	log.Println("In HandleFieldList")
	return h.remoteDb.FieldList(table, fieldWildcard)
}

// HandleStmtPrepare is called for COM_STMT_PREPARE
func (h ProxyRequestHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	log.Println("In HandleStmtPrepare")
	stmt, err := h.remoteDb.Prepare(query)
	if err != nil {
		return 0, 0, nil, err
	}
	return stmt.ParamNum(), stmt.ColumnNum(), stmt, nil
}

// 'context' isn't used but replacing it with `_` would remove important information for who
// wants to extend this later.
//revive:disable:unused-parameter

// HandleStmtExecute is called for COM_STMT_EXECUTE
func (h ProxyRequestHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	log.Println("In HandleStmtExecute")
	return context.(*client.Stmt).Execute(args...)
}

// HandleStmtClose is called for COM_STMT_CLOSE
func (h ProxyRequestHandler) HandleStmtClose(context interface{}) error {
	log.Println("In HandleStmtClose")
	return context.(*client.Stmt).Close()
}

// HandleOtherCommand is called for commands not handled elsewhere
func (h ProxyRequestHandler) HandleOtherCommand(cmd byte, data []byte) error {
	log.Printf("Received: OtherCommand: cmd=%x, data=%x", cmd, data)
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported now", cmd),
	)
}
