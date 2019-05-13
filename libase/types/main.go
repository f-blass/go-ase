package types

import (
	"fmt"
	"reflect"
	"time"
)

//go:generate go run ./gen_types.go

// ASEType reflects the data types ASE supports.
type ASEType int

// Type retuns an ASEType based on the name.
func Type(name string) ASEType {
	t, ok := string2type[name]
	if !ok {
		return ILLEGAL
	}
	return t
}

// String implements the Stringer interface.
func (t ASEType) String() string {
	s, ok := type2string[t]
	if !ok {
		return ""
	}
	return s
}

// GoReflectType returns the reflect.Type of the Go type used for the
// ASEType.
func (t ASEType) GoReflectType() reflect.Type {
	t, ok := type2reflect[t]
	if !ok {
		return nil
	}
	return t
}

// GoType returns the Go type used for the ASEType.
func (t ASEType) GoType() interface{} {
	t, ok := type2interface[t]
	if !ok {
		return nil
	}
	return t
}

// FromGoType returns the most fitting ASEType for the Go type.
func FromGoType(value interface{}) (ASEType, error) {
	switch value.(type) {
	case int64:
		return BIGINT, nil
	case float64:
		return FLOAT, nil
	case bool:
		return BIT, nil
	case []byte:
		return BINARY, nil
	case string:
		return CHAR, nil
	case time.Time:
		return BIGDATETIME, nil
	default:
		return ILLEGAL, fmt.Errorf("Invalid type for ASE: %v", value)
	}
}
