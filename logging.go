// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                        /[logging.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// configurable logging settings
var (
	// LogFunc specifies the external logging function to
	// use instead of the logging done in this package.
	//
	// Calls to LogFunc are buffered for output in a separate
	// goroutine, so LogBufferSize setting still applies.
	//
	// The signature accepts multiple args, to match log.Print(), but when
	// called it will be passed a single string with the error/info message.
	//
	LogFunc func(args ...interface{})

	// LogFile specifies the log file for logging logError() and logInfo()
	// output. If you don't specify it, no writing to file will be done.
	//
	// If you've specified LogFunc, this setting is has no effect
	// since you're using an external logging function.
	//
	LogFile string

	// LogBufferSize specifies the number of messages buffered in logChan.
	// If you don't specify it, initLog will set it to 1024 by default.
	//
	// To disable log buffering, set it to -1. This can be useful if you want
	// to see log messages in-order with code execution, not after the fact.
	// But this will slow it down while it waits for log messages to be written.
	//
	LogBufferSize int
)

// logChan is the channel into which messages are sent.
var logChan chan string

// logError outputs an error message to the standard output
// and optionally to a log file specified by LogFile.
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

// logOuput outputs a log message immediately.
//
// If LogFunc is specified, it calls LogFunc and exits.
//
// Otherwise it prints the message to the standard output and
// optionally writes to the log file specified by LogFile.
//
func logOutput(msg string) {
	if LogFunc != nil {
		LogFunc(msg)
		return
	}
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
