// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                            /[log.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// pl is like fmt.Println, but returns no values. It is only used for debugging.
func pl(args ...interface{}) { fmt.Println(args...) }

var _ = pl

// LogPrint prints a logging message to the standard output.
//
// Prefixes each line in the message with a timestamp.
//
// This is a convenience function that provides default logging
// output for this package while developing or debugging.
//
// You can assign it to Sender.LogFunc and Receiver.LogFunc.
//
func LogPrint(args ...interface{}) {
	const printMsg = true
	const logFile = ""
	logEnter(printMsg, logFile, args...)
} //                                                                    LogPrint

// MakeLogFunc creates and returns a function to use for default logging.
//
// This is a convenience function that provides default logging
// output for this package, for use while developing or debugging.
//
// You can assign this function to Sender.LogFunc and Receiver.LogFunc:
// E.g. sender.LogFunc = MakeLogFunc(true, "udpt.log")
//
// printMsg determines if each message should be printed to standard output.
//
// logFile specifies the log file into which to append.
//
func MakeLogFunc(printMsg bool, logFile string) func(args ...interface{}) {
	return func(args ...interface{}) {
		logEnter(printMsg, logFile, args...)
	}
} //                                                                 MakeLogFunc

// -----------------------------------------------------------------------------

// logChanSize specifies the number of messages buffered in logChan.
//
// To disable log buffering, set it to -1. This can be useful if you want
// to see log messages in-order with code execution, not after the fact.
// But this will slow it down while it waits for log messages to be written.
//
const logChanSize = 1024

// logChan is the channel into which log messages are sent.
var logChan chan logEntry

// logEntry contains a message to be printed and/or written to a log file.
type logEntry struct {
	printMsg bool
	logFile  string
	msg      string
} //                                                                    logEntry

// Output immediately prints msg to standard output and if
// logFile is not blank, appends the message to logFile.
func (ob *logEntry) Output() {
	if ob.printMsg {
		fmt.Println(ob.msg)
	}
	path := ob.logFile
	if path == "" {
		return
	}
	const mode = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(path, mode, 0644) // -> (*os.File, error)
	if err != nil {
		fmt.Println("ERROR 0xE5CB4B Failed opening "+path+":", err)
		return
	}
	n, err := file.WriteString(ob.msg + "\n")
	if n == 0 || err != nil {
		fmt.Println("ERROR 0xE81F3D Failed writing "+path+":", err)
	}
	err = file.Close()
	if err != nil {
		fmt.Println("ERROR 0xE2EC72 Failed closing "+path+":", err)
	}
} //                                                                      Output

// logEnter enters a log message (built from args) in the log queue (logChan).
//
// If logChanSize is negative, outputs the message immediately.
//
// printMsg: determines if the message should be printed to standard output.
//
// logFile: specifies the log file into which to append.
//
func logEnter(printMsg bool, logFile string, args ...interface{}) {
	if logChanSize == 0 {
		logInit()
	}
	msg := logMakeMessage(args...)
	entry := logEntry{printMsg: printMsg, logFile: logFile, msg: msg}
	if logChan == nil {
		logInit()
	}
	if logChanSize > 0 {
		logChan <- entry
	} else {
		entry.Output()
	}
} //                                                                    logEnter

// logInit initializes the logging queue (logChan) and launches
// a goroutine that receives and outputs log messages.
func logInit() {
	n := logChanSize
	if n < 0 {
		n = 1
	}
	logChan = make(chan logEntry, n)
	go func() {
		for entry := range logChan {
			entry.Output()
		}
	}()
} //                                                                     logInit

// logMakeMessage creates a log message by joining args.
//
// Prefixes each line with the current date/time.
//
func logMakeMessage(args ...interface{}) string {
	var (
		ts  = time.Now().String()[:19]
		msg = joinArgs("", args...)
	)
	// prefix each line with a timestamp
	var lines = strings.Split(msg, "\n")
	for i, line := range lines {
		lines[i] = ts + " " + line
	}
	msg = strings.Join(lines, "\n")
	return msg
} //                                                              logMakeMessage

// end
