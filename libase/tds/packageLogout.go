// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import "fmt"

var _ Package = (*LogoutPackage)(nil)

type LogoutPackage struct {
	Options uint8
}

func (pkg *LogoutPackage) ReadFrom(ch BytesChannel) error {
	var err error
	pkg.Options, err = ch.Uint8()
	if err != nil {
		return ErrNotEnoughBytes
	}

	if pkg.Options != 0 {
		return fmt.Errorf("unhandled logout option %d", pkg.Options)
	}

	return nil
}

func (pkg LogoutPackage) WriteTo(ch BytesChannel) error {
	err := ch.WriteByte(byte(TDS_LOGOUT))
	if err != nil {
		return err
	}

	return ch.WriteUint8(pkg.Options)
}

func (pkg LogoutPackage) String() string {
	return fmt.Sprintf("%T(%d)", pkg, pkg.Options)
}
