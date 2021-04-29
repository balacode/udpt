// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                      /[log_error.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"fmt"
	"strings"
)

// logError outputs an error message to the standard
// output and to a log file specified by LogFile.
//
// Returns an error value initialized with the message.
func logError(id uint32, args ...interface{}) error {
	sID := "0x" + strings.TrimLeft(fmt.Sprintf("%08X", id), "0")
	msg := joinArgs("ERROR "+sID+":", args...)
	logInfo(msg)
	return fmt.Errorf(msg)
} //                                                                    logError

// joinArgs returns a string built from a list of variadic arguments 'args',
// with some minimal formatting rules described as follows:
//
// Inserts a space between each argument, unless the preceding argument
// ends with '(', or the current argument begins with ')' or ':'.
//
// If a string argument in 'args' begins with '^', then the '^' is removed
// and the argument's string is quoted in single quotes without escaping it.
//
// If a string argument in 'args' ends with '^', then the '^' is removed
// and the next argument is quoted in the same way.
func joinArgs(prefix string, args ...interface{}) string {
	var (
		quoteNext bool
		lastChar  byte
		retBuf    bytes.Buffer
	)
	ws := func(s string) {
		_, _ = retBuf.WriteString(s)
	}
	ws(prefix)
	for i, arg := range args {
		s := fmt.Sprint(arg)
		firstChar := byte(0)
		if len(s) > 0 {
			firstChar = s[0]
		}
		if i > 0 &&
			lastChar != '(' &&
			firstChar != ')' &&
			firstChar != ':' {
			ws(" ")
		}
		q := quoteNext
		if strings.HasPrefix(s, "^") {
			q = true
			s = s[1:]
		}
		quoteNext = strings.HasSuffix(s, "^")
		if quoteNext {
			s = s[:len(s)-1]
		}
		if q {
			ws("'")
		}
		ws(s)
		if q {
			ws("'")
		}
		lastChar = 0
		if len(s) > 0 {
			lastChar = s[len(s)-1]
		}
	}
	return retBuf.String()
} //                                                                    joinArgs

// end
