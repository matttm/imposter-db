package main

import (
	"fmt"
	"log"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
)

type ProxyRequestHandler struct {
	remoteDb *client.Conn
}

func NewRemoteHandler(c *client.Conn) ProxyRequestHandler {
	return ProxyRequestHandler{remoteDb: c}
}

// UseDB is called for COM_INIT_DB
func (h ProxyRequestHandler) UseDB(dbName string) error {
	return h.UseDB(dbName)
}

// HandleQuery is called for COM_QUERY
func (h ProxyRequestHandler) HandleQuery(query string) (*mysql.Result, error) {
	log.Printf("Received: Query: %s", query)

	// These two queries are implemented for minimal support for MySQL Shell
	if query == `SET NAMES 'utf8mb4';` {
		return nil, nil
	}
	if query == `select concat(@@version, ' ', @@version_comment)` {
		r, err := mysql.BuildSimpleResultset([]string{"concat(@@version, ' ', @@version_comment)"}, [][]interface{}{
			{"8.0.11"},
		}, false)
		if err != nil {
			return nil, err
		}
		return mysql.NewResult(r), nil
	}

	return nil, fmt.Errorf("not supported now")
}

// HandleFieldList is called for COM_FIELD_LIST packets
// Note that COM_FIELD_LIST has been deprecated since MySQL 5.7.11
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_field_list.html
func (h ProxyRequestHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	return h.HandleFieldList(table, fieldWildcard)
}

// HandleStmtPrepare is called for COM_STMT_PREPARE
func (h ProxyRequestHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	return h.HandleStmtPrepare(query)
}

// 'context' isn't used but replacing it with `_` would remove important information for who
// wants to extend this later.
//revive:disable:unused-parameter

// HandleStmtExecute is called for COM_STMT_EXECUTE
func (h ProxyRequestHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	return h.HandleStmtExecute(context, query, args)
}

// HandleStmtClose is called for COM_STMT_CLOSE
func (h ProxyRequestHandler) HandleStmtClose(context interface{}) error {
	return h.HandleStmtClose(context)
}

// HandleOtherCommand is called for commands not handled elsewhere
func (h ProxyRequestHandler) HandleOtherCommand(cmd byte, data []byte) error {
	log.Printf("Received: OtherCommand: cmd=%x, data=%x", cmd, data)
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported now", cmd),
	)
}
