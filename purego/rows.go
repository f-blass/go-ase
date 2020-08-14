package purego

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"

	"github.com/SAP/go-ase/libase/tds"
)

var (
	_ driver.Rows                           = (*Rows)(nil)
	_ driver.RowsNextResultSet              = (*Rows)(nil)
	_ driver.RowsColumnTypeLength           = (*Rows)(nil)
	_ driver.RowsColumnTypeDatabaseTypeName = (*Rows)(nil)
)

type Rows struct {
	Conn   *Conn
	RowFmt *tds.RowFmtPackage
}

func (rows Rows) Columns() []string {
	// TODO ignore hidden columns
	response := make([]string, len(rows.RowFmt.Fmts))

	for i, fieldFmt := range rows.RowFmt.Fmts {
		// TODO check if RowFmt is wide and contains column label,
		// catalgoue, schema, table
		response[i] = fieldFmt.Name()
	}

	return response
}

func (rows *Rows) Close() error {
	for {
		err := rows.NextResultSet()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("go-ase: error consuming result sets: %w", err)
		}
	}

	return nil
}

func (rows *Rows) Next(dst []driver.Value) error {
	pkg, err := rows.Conn.Channel.NextPackage(true)
	if err != nil {
		return fmt.Errorf("go-ase: error reading next row package: %w", err)
	}

	var row *tds.RowPackage
	switch typed := pkg.(type) {
	case *tds.RowPackage:
		row = typed
	case *tds.RowFmtPackage:
		// TODO: should next return io.EOF if the result set is
		// finished?
		rows.RowFmt = typed
		return rows.Next(dst)
	case *tds.DonePackage:
		return io.EOF
	default:
		return fmt.Errorf("go-ase: unhandled package type %T: %v", pkg, pkg)
	}

	// TODO handle hidden columns
	if len(dst) != len(row.DataFields) {
		return fmt.Errorf("go-ase: received %d destinations for %d data fields",
			len(dst), len(row.DataFields))
	}

	for i := range dst {
		// TODO correctly transform data
		dst[i] = row.DataFields[i].Value()
	}

	return nil
}

func (rows *Rows) HasNextResultSet() bool {
	// TODO this doesn't seem good.
	return rows.NextResultSet() != nil
}

func (rows *Rows) NextResultSet() error {
	// discard all RowPackage until either end of communication or next
	// RowFmtPackage
	for {
		pkg, err := rows.Conn.Channel.NextPackage(false)
		if err != nil {
			if errors.Is(err, tds.ErrNoPackageReady) {
				return io.EOF
			}
			return fmt.Errorf("go-ase: error reading next package: %w", err)
		}

		switch rowFmt := pkg.(type) {
		case *tds.RowFmtPackage:
			rows.RowFmt = rowFmt
			return nil
		case *tds.RowPackage:
			continue
		case *tds.DonePackage:
			return io.EOF
		default:
			return fmt.Errorf("go-ase: unhandled package type %T: %v", pkg, pkg)
		}
	}
}

func (rows *Rows) ColumnTypeLength(index int) (int64, bool) {
	return rows.RowFmt.Fmts[index].MaxLength(), true
}

func (rows *Rows) ColumnTypeDatabaseTypeName(index int) string {
	return string(rows.RowFmt.Fmts[index].DataType())
}
