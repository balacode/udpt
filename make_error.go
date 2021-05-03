// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                     /[make_error.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
)

// makeError returns a new error value from joining id and args.
// The ID is formatted as a 6-digit hex string. e.g. "0xE12345"
func makeError(id uint32, args ...interface{}) error {
	sID := fmt.Sprintf("ERROR 0x%06X:", id)
	msg := joinArgs(sID, args...)
	ret := fmt.Errorf(msg)
	return ret
} //                                                                   makeError

// end
