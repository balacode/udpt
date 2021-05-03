// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                   /[string_funcs.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"strings"
)

// getPart returns the substring between 'prefix' and 'suffix'.
//
// When the prefix is blank, returns the part from the beginning of 's'.
//
// When the suffix is blank, returns the part up to the end of 's'.
// I.e. if prefix and suffix are both blank, returns 's'.
//
// When either prefix or suffix is not found, returns a zero-length string.
//
func getPart(s, prefix, suffix string) string {
	at := strings.Index(s, prefix)
	if at == -1 {
		return ""
	}
	s = s[at+len(prefix):]
	if suffix == "" {
		return s
	}
	at = strings.Index(s, suffix)
	if at == -1 {
		return ""
	}
	return s[:at]
} //                                                                     getPart

// joinArgs joins args into a single string, with a space between arguments.
func joinArgs(tag string, args ...interface{}) string {
	ar := make([]string, len(args))
	for i, arg := range args {
		ar[i] = fmt.Sprint(arg)
	}
	ret := strings.TrimSpace(tag + " " + strings.Join(ar, " "))
	return ret
} //                                                                    joinArgs

// padf suffixes a string with spaces to make sure it is at least
// 'minLength' characters wide. I.e. the string is left-aligned.
//
// If the string is wider than 'minLength', returns the string as it is.
//
func padf(minLength int, format string, args ...interface{}) string {
	format = fmt.Sprintf(format, args...)
	if len(format) < minLength {
		return format + strings.Repeat(" ", minLength-len(format))
	}
	return format
} //                                                                        padf

// end
