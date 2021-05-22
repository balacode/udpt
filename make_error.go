// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                     /[make_error.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"regexp"
	"strings"
)

// makeError returns a new error value from joining id and args.
// The ID is formatted as a 6-digit hex string. e.g. "0xE12345"
func makeError(id uint32, args ...interface{}) error {
	rx := regexp.MustCompile(`ERROR 0x[0-9a-fA-F]*: `)
	msg := joinArgs("", args...)
	msg = string(rx.ReplaceAll([]byte(msg), []byte("")))
	msg = fmt.Sprintf("ERROR 0x%06X: ", id) + msg
	msg = strings.TrimSpace(msg)
	return fmt.Errorf(msg)
} //                                                                   makeError

// end
