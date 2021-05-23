// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                     /[make_error.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// makeError returns a new error value from joining 'id' and 'a'.
// The ID is formatted as a 6-digit hex string. e.g. "0xE12345"
func makeError(id uint32, a ...interface{}) error {
	rx := regexp.MustCompile(`ERROR 0x[0-9a-fA-F]*: `)
	m := joinArgs("", a...)
	m = string(rx.ReplaceAll([]byte(m), []byte("")))
	m = fmt.Sprintf("ERROR 0x%06X: ", id) + m
	m = strings.TrimSpace(m)
	return errors.New(m)
} //                                                                   makeError

// end
