package main

import (
	"fmt"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/server"
)

type ProxyRequestHandler struct {
	// //handle COM_INIT_DB command, you can check whether the dbName is valid, or other.
	// UseDB(dbName string) error
	// //handle COM_QUERY command, like SELECT, INSERT, UPDATE, etc...
	// //If Result has a Resultset (SELECT, SHOW, etc...), we will send this as the response, otherwise, we will send Result
	// HandleQuery(query string) (*Result, error)
	// //handle COM_FILED_LIST command
	// HandleFieldList(table string, fieldWildcard string) ([]*Field, error)
	// //handle COM_STMT_PREPARE, params is the param number for this statement, columns is the column number
	// //context will be used later for statement execute
	// HandleStmtPrepare(query string) (params int, columns int, context interface{}, err error)
	// //handle COM_STMT_EXECUTE, context is the previous one set in prepare
	// //query is the statement prepare query, and args is the params for this statement
	// HandleStmtExecute(context interface{}, query string, args []interface{}) (*Result, error)
	// //handle COM_STMT_CLOSE, context is the previous one set in prepare
	// //this handler has no response
	// HandleStmtClose(context interface{}) error
	// //handle any other command that is not currently handled by the library,
	// //default implementation for this method will return an ER_UNKNOWN_ERROR
	// HandleOtherCommand(cmd byte, data []byte) errorHandler struct {
	handler server.Handler
}

// UseDB is called for COM_INIT_DB
func (h *ProxyRequestHandler) UseDB(dbName string) error {
	return h.UseDB(dbName)
}

// HandleQuery is called for COM_QUERY
func (h *ProxyRequestHandler) HandleQuery(query string) (*Result, error) {
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
func (h ProxyRequestHandler) HandleFieldList(table string, fieldWildcard string) ([]*Field, error) {
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
func (h ProxyRequestHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*Result, error) {
	return h.HandleStmtExecute(context, query, args)
}

// HandleStmtClose is called for COM_STMT_CLOSE
func (h ProxyRequestHandler) HandleStmtClose(context interface{}) error {
	return h.HandleStmtClose(context)
}
