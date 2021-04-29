// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[log_info.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	// LogFile specifies the name of the output log file.
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

// initLog initializes logging
func initLog() {
	if LogFile == "" {
		LogFile = os.Args[0] + ".log"
	}
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

// logInfo logs a message to the standard output and a log file
func logInfo(args ...interface{}) {
	if LogFile == "" || LogBufferSize == 0 {
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

// logOutput prints a message to the standard output
// and writes to the log file immediately
func logOutput(msg string) {
	fmt.Println(msg)
	file, err := os.OpenFile( // -> (*os.File, error)
		LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		file.WriteString(msg + "\n")
		n, err := file.WriteString(msg + "\n")
		if n == 0 || err != nil {
			_ = logError(0xE81F3D, "Failed writing", LogFile, ":", err)
		}
		err = file.Close()
		if err != nil {
			_ = logError(0xE2EC72, "Failed closing", LogFile, ":", err)
		}
	} else {
		_ = logError(0xE5CB4B, "Opening file", LogFile, ":", err)
	}
} //                                                                   logOutput

// end
