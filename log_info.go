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
	// LogFile _ _
	LogFile string

	// LogBufferSize  _ _
	LogBufferSize int

	// logChan _ _
	logChan chan string
)

// initLog _ _
func initLog() {
	if LogFile == "" {
		LogFile = os.Args[0] + ".log"
	}
	if LogBufferSize == 0 {
		LogBufferSize = 1024
	}
	logChan = make(chan string, LogBufferSize)
	go func() {
		for msg := range logChan {
			fmt.Println(msg)
			file, err := os.OpenFile( // -> (*os.File, error)
				LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err == nil {
				file.WriteString(msg + "\n")
				file.Close()
			} else {
				logError(0xE5CB4B, "Opening file", LogFile, ":", err)
			}
		}
	}()
} //                                                                     initLog

// logInfo _ _
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
	logChan <- msg
} //                                                                     logInfo

// end
