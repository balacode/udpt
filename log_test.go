// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[log_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run Test_log_*

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (le *logEntry) Output()
//
// go test -run Test_log_logEntry_Output_*

const logEntryTestFile = "udpt.logEntryTestFile.tmp"

// test writing to console without logging to a file
func Test_log_logEntry_Output_1(t *testing.T) {
	var tlog strings.Builder
	_ = os.Remove(logEntryTestFile)
	const s = "test entry #3146250897"
	le := logEntry{printMsg: true, logFile: "", msg: s}
	//
	le.outputDI(&tlog, nil)
	//
	ts := tlog.String()
	if ts != s+"\n" {
		t.Error("0xE58FB6")
	}
	_ = os.Remove(logEntryTestFile)
}

// test logging to a file without writing to console
func Test_log_logEntry_Output_2(t *testing.T) {
	var tlog strings.Builder
	_ = os.Remove(logEntryTestFile)
	const s = "test entry #9473258061"
	le := logEntry{printMsg: false, logFile: logEntryTestFile, msg: s}
	//
	le.outputDI(&tlog, openLogFile)
	//
	ts := tlog.String()
	if ts != "" {
		t.Error("0xE8BB8F")
	}
	data, err := ioutil.ReadFile(logEntryTestFile)
	if string(data) != s+"\n" || err != nil {
		t.Error("0xED98CA")
	}
	_ = os.Remove(logEntryTestFile)
}

// test that no file is created when openLogFile() returns nil
func Test_log_logEntry_Output_3(t *testing.T) {
	var tlog strings.Builder
	_ = os.Remove(logEntryTestFile)
	const s = "test entry #2361498057"
	openLogFile := func(string) (io.WriteCloser, error) {
		return nil, makeError(0xEE05A5, "failed openLogFile")
	}
	le := logEntry{printMsg: false, logFile: logEntryTestFile, msg: s}
	//
	le.outputDI(&tlog, openLogFile)
	//
	data, err := ioutil.ReadFile(logEntryTestFile)
	if data != nil || !os.IsNotExist(err) {
		t.Error("0xE7BC99")
	}
	_ = os.Remove(logEntryTestFile)
}

// test that errors are written to console when Write() and Close() fail
func Test_log_logEntry_Output_4(t *testing.T) {
	_ = os.Remove(logEntryTestFile)
	const s = "test entry #1540938672"
	openLogFile := func(string) (io.WriteCloser, error) {
		return &mockWriteCloser{failWrite: true, failClose: true}, nil
	}
	var tlog strings.Builder
	le := logEntry{printMsg: false, logFile: logEntryTestFile, msg: s}
	le.outputDI(&tlog, openLogFile)
	//
	// log must contain two error descriptions
	ts := tlog.String()
	if !strings.Contains(ts, "ERROR 0x") ||
		!strings.Contains(ts, "failed mockWriteCloser.Write") ||
		!strings.Contains(ts, "failed mockWriteCloser.Close") {
		t.Error("0xEA1B28")
	}
	// must not create/write to file
	data, err := ioutil.ReadFile(logEntryTestFile)
	if data != nil || !os.IsNotExist(err) {
		t.Error("0xEB38C2")
	}
	_ = os.Remove(logEntryTestFile)
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// openLogFile(filename string, logErrorTo io.Writer) io.WriteCloser
//
// go test -run Test_log_openLogFile_*

// must succeed opening a file and writing to it two times
func Test_log_openLogFile_1(t *testing.T) {
	//
	const testfile = "udpt.Test_log_openLogFile_.tmp"
	_ = os.Remove(testfile)
	var tlog strings.Builder
	now := time.Now().String()
	msg1 := now + " msg1\n"
	msg2 := now + " msg2\n"
	//
	wrc, err := openLogFile(testfile)
	if err != nil {
		t.Error("0xEA03F4", err)
	}
	wrc.Write([]byte(msg1))
	wrc.Close()
	//
	wrc, err = openLogFile(testfile)
	if err != nil {
		t.Error("0xED8C37", err)
	}
	wrc.Write([]byte(msg2))
	wrc.Close()
	//
	data, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Error("0xE7B84B", err)
	}
	content := string(data)
	if content != msg1+msg2 {
		t.Error("0xE24C06")
	}
	ts := tlog.String()
	if ts != "" {
		t.Error("0xE4CE4D")
	}
	_ = os.Remove(testfile)
}

// must fail to open a file with an invalid name
func Test_log_openLogFile_2(t *testing.T) {
	wrc, err := openLogFile("\\:/")
	if wrc != nil {
		t.Error("0xEC8BA1")
	}
	if !matchError(err, "ERROR 0x") {
		t.Error("0xEE9DE4", "wrong error:", err)
	}
}

// -----------------------------------------------------------------------------

// (lw *logWriter) initChan()
//
// go test -run Test_log_initChan_
//
func Test_log_initChan_(t *testing.T) {
	lw := &logWriter{}
	lw.initChan()
	if lw.logChan == nil {
		t.Error("0xE39B91", "logChan not initialized")
	}
}

// end
