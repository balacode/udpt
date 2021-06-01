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

// logChanSize specifies the number of log messages buffered in logChan.
//
// To disable log buffering, set it to 1. This can be useful if you want
// to see log messages in-order with code execution, not after the fact.
// But this will slow it down while it waits for log messages to be written.
//
const logChanSize = 1024

// logTimeNow should always be set to time.Now,
// except when mocking during testing.
var logTimeNow = time.Now

// logWriter provides a buffered logger that outputs
// to standard output and/or a log file.
type logWriter struct {
	Print   bool
	LogFile string

	// logChan is the channel into which log messages are sent.
	logChan chan logEntry
}

// MakeLogWriter is a convenience function that creates and returns
// a Writer to use for logging during developing or debugging.
//
// You can assign this function to Config.LogFunc:
// E.g. sender.Config.LogWriter = MakeLogWriter(true, "udpt.log")
//
// printMsg determines if each log entry should be printed to standard output.
//
// logFile specifies the log file into which to append.
//
func MakeLogWriter(printMsg bool, logFile string) io.Writer {
	return &logWriter{Print: printMsg, LogFile: logFile}
} //                                                               MakeLogWriter

// Write writes 'b' to the log and implements io.Writer
func (lw *logWriter) Write(b []byte) (n int, err error) {
	lw.logEnter(string(b))
	return len(b), nil
}

// initChan initializes the logging queue (logChan) and launches
// a goroutine that receives and outputs log messages.
func (lw *logWriter) initChan() {
	lw.logChan = make(chan logEntry, logChanSize)
	go func() {
		for entry := range lw.logChan {
			entry.Output()
		}
	}()
} //                                                                    initChan

// logEnter sends a log message to log channel 'logChan'.
// Prefixes each line in the message with a timestamp.
//
// If logChanSize is 1, outputs the message immediately.
func (lw *logWriter) logEnter(msg string) {
	//
	// prefix each line with a timestamp
	tms := logTimeNow().String()[:19]
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		lines[i] = tms + " " + line
	}
	msg = strings.Join(lines, "\n")
	entry := logEntry{
		printMsg: lw.Print,
		logFile:  lw.LogFile,
		msg:      msg,
	}
	if lw.logChan == nil {
		lw.initChan()
	}
	lw.logChan <- entry
} //                                                                    logEnter

// -----------------------------------------------------------------------------

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
	dest io.Writer,
	openLogFile func(filename string) (io.WriteCloser, error),
) {
	if le.printMsg {
		fmt.Fprintln(dest, le.msg)
	}
	path := le.logFile
	if path == "" {
		return
	}
	wr, err := openLogFile(path)
	if err != nil {
		fmt.Fprintln(dest, err)
		return
	}
	if wr == nil {
		return
	}
	n, err := wr.Write([]byte(le.msg + "\n"))
	if n == 0 || err != nil {
		fmt.Fprintln(dest, "ERROR 0xE81F3D:", err)
	}
	err = wr.Close()
	if err != nil {
		fmt.Fprintln(dest, "ERROR 0xE50F96:", err)
	}
} //                                                                    outputDI

// openLogFile opens a file for outputDI().
func openLogFile(filename string) (io.WriteCloser, error) {
	fl, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, makeError(0xE2DA59, err)
	}
	return fl, nil
} //                                                                 openLogFile

// end
