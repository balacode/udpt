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

// -----------------------------------------------------------------------------

// pl(a ...interface{})
//
// go test -run Test_log_pl_
//
func Test_log_pl_(t *testing.T) {
	pl()
} //                                                                Test_log_pl_

// LogPrint(a ...interface{})
//
// go test -run Test_log_LogPrint_
//
func Test_log_LogPrint_(t *testing.T) {
	LogPrint()
} //                                                          Test_log_LogPrint_

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// MakeLogFunc(printMsg bool, logFile string) func(a ...interface{})
//
// go test -run Test_log_MakeLogFunc_*

func Test_log_MakeLogFunc_1(t *testing.T) {
	//
	// prepare: delete log file and mock time.Now()
	logFile := "udpt.Test_log_MakeLogFunc_.tmp"
	_ = os.Remove(logFile)
	var tm = time.Now()
	logTimeNow = func() time.Time {
		return tm
	}
	// test!
	fn := MakeLogFunc(false, logFile)
	if fn == nil {
		t.Error("0xEA9A06", "MakeLogFunc returned nil")
	}
	fn("abc", 123, "45", "de")
	//
	// check results
	time.Sleep(500 * time.Millisecond) // wait for logger finish writing
	data, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Error("0xE08D49", err)
		return
	}
	got := string(data)
	expect := "@tm abc 123 45 de\n"
	expect = strings.ReplaceAll(expect, "@tm", tm.String()[:19])
	//
	if got != expect {
		t.Errorf("0xEF57A5"+" wrong text in log file "+
			"\n expect: %#v"+
			"\n    got: %#v",
			expect, got)
	}
	// cleanup
	_ = os.Remove(logFile)
	logTimeNow = time.Now
} //                                                      Test_log_MakeLogFunc_1

func Test_log_MakeLogFunc_2(t *testing.T) {
	const testFile = "udpt.Test_log_MakeLogFunc_.tmp"
	const msg = "test message #5067124389"
	_ = os.Remove(testFile)
	fn := MakeLogFunc(false, testFile)
	fn(msg)
	time.Sleep(time.Second)
	data, err := ioutil.ReadFile(testFile)
	if !strings.Contains(string(data), msg) || err != nil {
		t.Error("0xEC7ED0")
	}
	_ = os.Remove(testFile)
} //                                                      Test_log_MakeLogFunc_2

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (le *logEntry) Output()
//
// go test -run Test_log_logEntry_Output_*

const logEntryTestFile = "udpt.logEntryTestFile.tmp"

func Test_log_logEntry_Output_1(t *testing.T) {
	_ = os.Remove(logEntryTestFile)
	//
	// test writing to console without logging to a file
	const msg = "test message #3146250897"
	l := logEntry{printMsg: true, logFile: "", msg: msg}
	var sb strings.Builder
	l.outputDI(&sb, nil)
	if sb.String() != msg+"\n" {
		t.Error("0xE58FB6")
	}
	_ = os.Remove(logEntryTestFile)
} //                                                  Test_log_logEntry_Output_1

func Test_log_logEntry_Output_2(t *testing.T) {
	_ = os.Remove(logEntryTestFile)
	//
	// test logging to a file without writing to console
	const msg = "test message #9473258061"
	l := logEntry{printMsg: false, logFile: logEntryTestFile, msg: msg}
	var sb strings.Builder
	l.outputDI(&sb, openLogFile)
	if sb.String() != "" {
		t.Error("0xE8BB8F")
	}
	data, err := ioutil.ReadFile(logEntryTestFile)
	if string(data) != msg+"\n" || err != nil {
		t.Error("0xED98CA")
	}
	_ = os.Remove(logEntryTestFile)
} //                                                  Test_log_logEntry_Output_2

func Test_log_logEntry_Output_3(t *testing.T) {
	_ = os.Remove(logEntryTestFile)
	//
	// test that no file is created when openLogFile() returns nil
	const msg = "test message #2361498057"
	l := logEntry{printMsg: false, logFile: logEntryTestFile, msg: msg}
	openLogFile := func(string, io.Writer) io.WriteCloser { return nil }
	var sb strings.Builder
	l.outputDI(&sb, openLogFile)
	data, err := ioutil.ReadFile(logEntryTestFile)
	if data != nil || !os.IsNotExist(err) {
		t.Error("0xE7BC99")
	}
	_ = os.Remove(logEntryTestFile)
} //                                                  Test_log_logEntry_Output_3

func Test_log_logEntry_Output_4(t *testing.T) {
	_ = os.Remove(logEntryTestFile)
	//
	// test that errors are written to console when Write() and Close() fail
	const msg = "test message #1540938672"
	l := logEntry{printMsg: false, logFile: logEntryTestFile, msg: msg}
	openLogFile := func(string, io.Writer) io.WriteCloser {
		return &mockWriteCloser{failWrite: true, failClose: true}
	}
	var sb strings.Builder
	l.outputDI(&sb, openLogFile)
	//
	// console must contain two error messages
	cons := sb.String()
	if !strings.Contains(cons, "ERROR 0x") ||
		!strings.Contains(cons, "failed mockWriteCloser.Write") ||
		!strings.Contains(cons, "failed mockWriteCloser.Close") {
		t.Error("0xEA1B28")
	}
	// must not create/write to file
	data, err := ioutil.ReadFile(logEntryTestFile)
	if data != nil || !os.IsNotExist(err) {
		t.Error("0xEB38C2")
	}
	_ = os.Remove(logEntryTestFile)
} //                                                  Test_log_logEntry_Output_4

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// openLogFile(filename string, cons io.Writer) io.WriteCloser

// go test -run Test_log_openLogFile_1
//
func Test_log_openLogFile_1(t *testing.T) {
	//
	// must succeed opening a file and writing to it two times
	const testfile = "udpt.Test_log_openLogFile_.tmp"
	_ = os.Remove(testfile)
	var sb strings.Builder
	now := time.Now().String()
	msg1 := now + " msg1\n"
	msg2 := now + " msg2\n"
	//
	wrc := openLogFile(testfile, &sb)
	wrc.Write([]byte(msg1))
	wrc.Close()
	//
	wrc = openLogFile(testfile, &sb)
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
	if sb.String() != "" {
		t.Error("0xE4CE4D")
	}
	_ = os.Remove(testfile)
} //                                                      Test_log_openLogFile_1

// go test -run Test_log_openLogFile_2
//
func Test_log_openLogFile_2(t *testing.T) {
	//
	// must fail to open a file with an invalid name
	var sb strings.Builder
	wrc := openLogFile("\\:/", &sb)
	if wrc != nil {
		t.Error("0xEC8BA1")
	}
	if !strings.Contains(sb.String(), "ERROR 0x") {
		t.Error("0xEA1DC3")
	}
} //                                                      Test_log_openLogFile_2

// -----------------------------------------------------------------------------

// logInit()
//
// go test -run Test_log_logInit_
//
func Test_log_logInit_(t *testing.T) {
	logInit()
	if logChan == nil {
		t.Error("0xE39B91", "logChan not initialized")
	}
} //                                                           Test_log_logInit_

// logMakeMessage(tm time.Time, a ...interface{}) string
//
// go test -run Test_log_logMakeMessage_
//
func Test_log_logMakeMessage_(t *testing.T) {
	for _, it := range []struct {
		a      []interface{}
		expect string
	}{
		{nil, "@tm "},
		{[]interface{}{"line1\nline2"}, "@tm line1\n@tm line2"},
		{[]interface{}{"abc", 123, "de", 34}, "@tm abc 123 de 34"},
	} {
		var (
			tm     = time.Now()
			tms    = tm.String()[:19]
			got    = logMakeMessage(tm, it.a...)
			expect = strings.ReplaceAll(it.expect, "@tm", tms)
		)
		if got != expect {
			t.Errorf("0xEB14AF"+" logMakeMessage(<%s>, %#v)"+
				"\n expect: %#v"+
				"\n    got: %#v",
				tms, it.a, expect, got)
		}
	}
} //                                                    Test_log_logMakeMessage_

// end
