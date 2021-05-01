// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                        /[logging.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	// LogFile specifies the log file for logging logError() and logInfo()
	// output. If you don't specify it, no writing to file will be done.
	LogFile string

	// LogBufferSize specifies the number of messages buffered in logChan.
	// If you don't specify it, initLog will set it to 1024 by default.
	//
	// To disable log buffering, set it to -1. This can be useful if you want
	// to see log messages in-order with code execution, not after the fact.
	// But this will slow it down while it waits for log messages to be written.
	LogBufferSize int

	// logChan is the channel into which messages are sent
	logChan chan string
)

// logError outputs an error message to the standard
// output and to a log file specified by LogFile.
//
// Returns an error value initialized with the message.
//
func logError(id uint32, args ...interface{}) error {
	sID := "0x" + strings.TrimLeft(fmt.Sprintf("%08X", id), "0")
	msg := joinArgs("ERROR "+sID+":", args...)
	logInfo(msg)
	return fmt.Errorf(msg)
} //                                                                    logError

// logInfo logs a message to the standard output and
// optionally to the log file specified by LogFile.
func logInfo(args ...interface{}) {
	if LogBufferSize == 0 {
		initLog()
	}
	ts := time.Now().String()[:19]
	//
	// join all the parts into a single string
	var msg string
	{
		strs := make([]string, len(args))
		for i, arg := range args {
			strs[i] = fmt.Sprint(arg)
		}
		msg = strings.Join(strs, " ")
	}
	// prefix each line with a timestamp
	var lines = strings.Split(msg, "\n")
	for i, line := range lines {
		lines[i] = ts + " " + line
	}
	msg = strings.Join(lines, "\n")
	if LogBufferSize > 0 {
		logChan <- msg
	} else {
		logOutput(msg)
	}
} //                                                                     logInfo

// -----------------------------------------------------------------------------
// # Helper Functions

// initLog initializes logging.
func initLog() {
	if LogBufferSize == 0 {
		LogBufferSize = 1024
	}
	size := LogBufferSize
	if size < 0 {
		size = 1
	}
	logChan = make(chan string, size)
	go func() {
		for msg := range logChan {
			logOutput(msg)
		}
	}()
} //                                                                     initLog

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
//
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

// logOutput prints a message to the standard output
// and writes to the log file immediately
func logOutput(msg string) {
	fmt.Println(msg)
	if LogFile == "" {
		return
	}
	const mode = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(LogFile, mode, 0644) // -> (*os.File, error)
	if err != nil {
		_ = logError(0xE5CB4B, "Failed opening", LogFile, ":", err)
		return
	}
	n, err := file.WriteString(msg + "\n")
	if n == 0 || err != nil {
		_ = logError(0xE81F3D, "Failed writing", LogFile, ":", err)
	}
	err = file.Close()
	if err != nil {
		_ = logError(0xE2EC72, "Failed closing", LogFile, ":", err)
	}
} //                                                                   logOutput

// end
