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
		for s := range logChan {
			fmt.Println(s)
			file, err := os.OpenFile( // -> (*os.File, error)
				LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err == nil {
				file.WriteString(s + "\n")
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
	var s string
	{
		strs := make([]string, len(args))
		for i, arg := range args {
			strs[i] = fmt.Sprint(arg)
		}
		s = strings.Join(strs, " ")
	}
	// prefix each line with a timestamp
	var lines = strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = ts + " " + line
	}
	s = strings.Join(lines, "\n")
	logChan <- s
} //                                                                     logInfo

// end
