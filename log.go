// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                            /[log.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// logTimeNow should always be set to time.Now,
// except when mocking during testing.
var logTimeNow = time.Now

// LogPrint prints a logging message to the standard output.
//
// Prefixes each line in the message with a timestamp.
//
// This is a convenience function that provides default logging
// output for this package while developing or debugging.
//
// You can assign it to Config.LogFunc
//
func LogPrint(a ...interface{}) {
	const printMsg = true
	const logFile = ""
	logEnter(printMsg, logFile, a...)
} //                                                                    LogPrint

// MakeLogFunc creates and returns a function to use for default logging.
//
// This is a convenience function that provides default logging
// output for this package, for use while developing or debugging.
//
// You can assign this function to Config.LogFunc:
// E.g. sender.LogFunc = MakeLogFunc(true, "udpt.log")
//
// printMsg determines if each message should be printed to standard output.
//
// logFile specifies the log file into which to append.
//
func MakeLogFunc(printMsg bool, logFile string) func(a ...interface{}) {
	return func(a ...interface{}) {
		logEnter(printMsg, logFile, a...)
	}
} //                                                                 MakeLogFunc

// -----------------------------------------------------------------------------

// logChanSize specifies the number of messages buffered in logChan.
//
// To disable log buffering, set it to 1. This can be useful if you want
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
func (le *logEntry) Output() {
	le.outputDI(os.Stdout, openLogFile)
} //                                                                      Output

// outputDI is only used by Output() and provides parameters
// for dependency injection, to enable mocking during testing.
func (le *logEntry) outputDI(
	cons io.Writer,
	openLogFile func(filename string, cons io.Writer) io.WriteCloser,
) {
	if le.printMsg {
		fmt.Fprintln(cons, le.msg)
	}
	path := le.logFile
	if path == "" {
		return
	}
	wr := openLogFile(path, cons)
	if wr == nil {
		return
	}
	n, err := wr.Write([]byte(le.msg + "\n"))
	if n == 0 || err != nil {
		fmt.Fprintln(cons, "ERROR 0xE81F3D:", err)
	}
	err = wr.Close()
	if err != nil {
		fmt.Fprintln(cons, "ERROR 0xE50F96:", err)
	}
} //                                                                    outputDI

// openLogFile opens a file for outputDI().
func openLogFile(filename string, cons io.Writer) io.WriteCloser {
	fl, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(cons, "ERROR 0xE2DA59:", err)
		return nil
	}
	return fl
} //                                                                 openLogFile

// -----------------------------------------------------------------------------

// logEnter sends a log message built from args 'a' to log channel 'logChan'.
//
// If logChanSize is 1, outputs the message immediately.
//
// printMsg: determines if the message should be printed to standard output.
//
// logFile: specifies the log file into which to append.
//
func logEnter(printMsg bool, logFile string, a ...interface{}) {
	msg := logMakeMessage(logTimeNow(), a...)
	entry := logEntry{printMsg: printMsg, logFile: logFile, msg: msg}
	if logChan == nil {
		logInit()
	}
	logChan <- entry
} //                                                                    logEnter

// logInit initializes the logging queue (logChan) and launches
// a goroutine that receives and outputs log messages.
func logInit() {
	logChan = make(chan logEntry, logChanSize)
	go func() {
		for entry := range logChan {
			entry.Output()
		}
	}()
} //                                                                     logInit

// logMakeMessage creates a log message by joining arguments in 'a'.
//
// Prefixes each line with the date/time specified in 'tm'.
//
func logMakeMessage(tm time.Time, a ...interface{}) string {
	var (
		tms = tm.String()[:19]
		msg = joinArgs("", a...)
	)
	// prefix each line with a timestamp
	var lines = strings.Split(msg, "\n")
	for i, line := range lines {
		lines[i] = tms + " " + line
	}
	msg = strings.Join(lines, "\n")
	return msg
} //                                                              logMakeMessage

// end
