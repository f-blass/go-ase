// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package cgo

//#include <stdlib.h>
//#include "ctlib.h"
import "C"
import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"unsafe"

	"github.com/SAP/go-ase/libase"
	"github.com/SAP/go-ase/libase/libdsn"
	"github.com/SAP/go-ase/libase/types"
)

// connection is the struct which represents a database connection.
type Connection struct {
	conn      *C.CS_CONNECTION
	driverCtx *csContext
}

// Interface satisfaction checks
var (
	_ driver.Conn               = (*Connection)(nil)
	_ driver.ConnBeginTx        = (*Connection)(nil)
	_ driver.ConnPrepareContext = (*Connection)(nil)
	_ driver.Execer             = (*Connection)(nil)
	_ driver.ExecerContext      = (*Connection)(nil)
	_ driver.Pinger             = (*Connection)(nil)
	_ driver.Queryer            = (*Connection)(nil)
	_ driver.QueryerContext     = (*Connection)(nil)
	_ driver.NamedValueChecker  = (*Connection)(nil)
)

// newConnection allocated initializes a new connection based on the
// options in the dsn.
//
// If driverCtx is nil a new csContext will be initialized.
func NewConnection(driverCtx *csContext, dsn *libdsn.Info) (*Connection, error) {
	if driverCtx == nil {
		var err error
		driverCtx, err = newCsContext(dsn)
		if err != nil {
			return nil, fmt.Errorf("Failed to initialize context for conn: %w", err)
		}
	}

	err := driverCtx.newConn()
	if err != nil {
		return nil, fmt.Errorf("Failed to ensure context: %w", err)
	}

	conn := &Connection{
		driverCtx: driverCtx,
	}

	retval := C.ct_con_alloc(driverCtx.ctx, &conn.conn)
	if retval != C.CS_SUCCEED {
		conn.Close()
		return nil, makeError(retval, "C.ct_con_alloc failed")
	}

	// Set password encryption
	cTrue := C.CS_TRUE
	retval = C.ct_con_props(conn.conn, C.CS_SET, C.CS_SEC_EXTENDED_ENCRYPTION,
		unsafe.Pointer(&cTrue), C.CS_UNUSED, nil)
	if retval != C.CS_SUCCEED {
		conn.Close()
		return nil, makeError(retval, "C.ct_con_props failed for CS_SEC_EXTENDED_ENCRYPTION")
	}

	cFalse := C.CS_FALSE
	retval = C.ct_con_props(conn.conn, C.CS_SET, C.CS_SEC_NON_ENCRYPTION_RETRY,
		unsafe.Pointer(&cFalse), C.CS_UNUSED, nil)
	if retval != C.CS_SUCCEED {
		conn.Close()
		return nil, makeError(retval, "C.ct_con_props failed for CS_SEC_NON_ENCRYPTION_RETRY")
	}

	// Give preference to the user store key
	if len(dsn.Userstorekey) > 0 {
		// Set userstorekey
		userstorekey := unsafe.Pointer(C.CString(dsn.Userstorekey))
		defer C.free(userstorekey)
		retval = C.ct_con_props(conn.conn, C.CS_SET, C.CS_SECSTOREKEY, userstorekey, C.CS_NULLTERM, nil)
		if retval != C.CS_SUCCEED {
			conn.Close()
			return nil, makeError(retval, "C.ct_con_props failed for C.CS_SECSTOREKEY")
		}
	} else {
		// Set username.
		username := unsafe.Pointer(C.CString(dsn.Username))
		defer C.free(username)
		retval = C.ct_con_props(conn.conn, C.CS_SET, C.CS_USERNAME, username, C.CS_NULLTERM, nil)
		if retval != C.CS_SUCCEED {
			conn.Close()
			return nil, makeError(retval, "C.ct_con_props failed for CS_USERNAME")
		}

		// Set password.
		password := unsafe.Pointer(C.CString(dsn.Password))
		defer C.free(password)
		retval = C.ct_con_props(conn.conn, C.CS_SET, C.CS_PASSWORD, password, C.CS_NULLTERM, nil)
		if retval != C.CS_SUCCEED {
			conn.Close()
			return nil, makeError(retval, "C.ct_con_props failed for CS_PASSWORD")
		}
	}

	if dsn.Host != "" && dsn.Port != "" {
		// Set hostname and port as string, since it is modified if
		// '-o ssl' is set.
		strHostport := dsn.Host + " " + dsn.Port
		// If '-o ssl='-option is set, add it to strHostport
		if dsn.TLSHostname != "" {
			strHostport += fmt.Sprintf("ssl=\"%s\"", dsn.TLSHostname)
		}
		// Create pointer
		ptrHostport := unsafe.Pointer(C.CString(strHostport))
		defer C.free(ptrHostport)

		retval = C.ct_con_props(conn.conn, C.CS_SET, C.CS_SERVERADDR, ptrHostport, C.CS_NULLTERM, nil)
		if retval != C.CS_SUCCEED {
			conn.Close()
			return nil, makeError(retval, "C.ct_con_props failed for CS_SERVERADDR")
		}
	}

	retval = C.ct_connect(conn.conn, nil, 0)
	if retval != C.CS_SUCCEED {
		conn.Close()
		return nil, makeError(retval, "C.ct_connect failed")
	}

	// Set database
	if dsn.Database != "" {
		_, err := conn.Exec("use "+dsn.Database, nil)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("Failed to connect to database %s: %w", dsn.Database, err)
		}
	}

	return conn, nil
}

// Close closes and deallocates a connection.
func (conn *Connection) Close() error {
	// Call context.drop when exiting this function to decrease the
	// connection counter and potentially deallocate the context.
	defer conn.driverCtx.dropConn()

	retval := C.ct_close(conn.conn, C.CS_UNUSED)
	if retval != C.CS_SUCCEED {
		return makeError(retval, "C.ct_close failed, connection has results pending")
	}

	retval = C.ct_con_drop(conn.conn)
	if retval != C.CS_SUCCEED {
		return makeError(retval, "C.ct_con_drop failed")
	}

	conn.conn = nil
	return nil
}

func (conn *Connection) Ping(ctx context.Context) error {
	rows, err := conn.QueryContext(ctx, "SELECT 'PING'", nil)
	if err != nil {
		return driver.ErrBadConn
	}
	defer rows.Close()

	cols := rows.Columns()
	cellRefs := make([]driver.Value, len(cols))

	for {
		err := rows.Next(cellRefs)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error occurred while exhausting result set: %w", err)
		}
	}

	return nil
}

func (conn *Connection) Exec(query string, args []driver.Value) (driver.Result, error) {
	return conn.ExecContext(context.Background(), query, libase.ValuesToNamedValues(args))
}

func (conn *Connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	_, result, err := conn.GenericExec(ctx, query, args)
	return result, err
}

func (conn *Connection) Query(query string, args []driver.Value) (driver.Rows, error) {
	return conn.QueryContext(context.Background(), query, libase.ValuesToNamedValues(args))
}

func (conn *Connection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	rows, _, err := conn.GenericExec(ctx, query, args)
	return rows, err
}

func (conn *Connection) CheckNamedValue(nv *driver.NamedValue) error {
	v, err := types.DefaultValueConverter.ConvertValue(nv.Value)
	if err != nil {
		return err
	}

	nv.Value = v
	return nil
}
