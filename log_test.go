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
}

// LogPrint(a ...interface{})
//
// go test -run Test_log_LogPrint_
//
func Test_log_LogPrint_(t *testing.T) {
	LogPrint()
}

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
	want := "@tm abc 123 45 de\n"
	want = strings.ReplaceAll(want, "@tm", tm.String()[:19])
	//
	if got != want {
		t.Errorf("0xEF57A5"+" wrong text in log file "+
			"\n want: %#v"+
			"\n  got: %#v",
			want, got)
	}
	// cleanup
	_ = os.Remove(logFile)
	logTimeNow = time.Now
}

func Test_log_MakeLogFunc_2(t *testing.T) {
	const testFile = "udpt.Test_log_MakeLogFunc_.tmp"
	const s = "test message #5067124389"
	_ = os.Remove(testFile)
	fn := MakeLogFunc(false, testFile)
	fn(s)
	time.Sleep(time.Second)
	data, err := ioutil.ReadFile(testFile)
	if !strings.Contains(string(data), s) || err != nil {
		t.Error("0xEC7ED0")
	}
	_ = os.Remove(testFile)
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (le *logEntry) Output()
//
// go test -run Test_log_logEntry_Output_*

const logEntryTestFile = "udpt.logEntryTestFile.tmp"

func Test_log_logEntry_Output_1(t *testing.T) {
	_ = os.Remove(logEntryTestFile)
	//
	// test writing to console without logging to a file
	const s = "test message #3146250897"
	le := logEntry{printMsg: true, logFile: "", msg: s}
	var tlog strings.Builder
	le.outputDI(&tlog, nil)
	if tlog.String() != s+"\n" {
		t.Error("0xE58FB6")
	}
	_ = os.Remove(logEntryTestFile)
}

func Test_log_logEntry_Output_2(t *testing.T) {
	_ = os.Remove(logEntryTestFile)
	//
	// test logging to a file without writing to console
	const s = "test message #9473258061"
	le := logEntry{printMsg: false, logFile: logEntryTestFile, msg: s}
	var tlog strings.Builder
	le.outputDI(&tlog, openLogFile)
	if tlog.String() != "" {
		t.Error("0xE8BB8F")
	}
	data, err := ioutil.ReadFile(logEntryTestFile)
	if string(data) != s+"\n" || err != nil {
		t.Error("0xED98CA")
	}
	_ = os.Remove(logEntryTestFile)
}

func Test_log_logEntry_Output_3(t *testing.T) {
	_ = os.Remove(logEntryTestFile)
	//
	// test that no file is created when openLogFile() returns nil
	const s = "test message #2361498057"
	le := logEntry{printMsg: false, logFile: logEntryTestFile, msg: s}
	openLogFile := func(string, io.Writer) io.WriteCloser { return nil }
	var tlog strings.Builder
	le.outputDI(&tlog, openLogFile)
	data, err := ioutil.ReadFile(logEntryTestFile)
	if data != nil || !os.IsNotExist(err) {
		t.Error("0xE7BC99")
	}
	_ = os.Remove(logEntryTestFile)
}

func Test_log_logEntry_Output_4(t *testing.T) {
	_ = os.Remove(logEntryTestFile)
	//
	// test that errors are written to console when Write() and Close() fail
	const s = "test message #1540938672"
	le := logEntry{printMsg: false, logFile: logEntryTestFile, msg: s}
	openLogFile := func(string, io.Writer) io.WriteCloser {
		return &mockWriteCloser{failWrite: true, failClose: true}
	}
	var tlog strings.Builder
	le.outputDI(&tlog, openLogFile)
	//
	// console must contain two error messages
	cons := tlog.String()
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
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// openLogFile(filename string, cons io.Writer) io.WriteCloser

// go test -run Test_log_openLogFile_1
//
func Test_log_openLogFile_1(t *testing.T) {
	//
	// must succeed opening a file and writing to it two times
	const testfile = "udpt.Test_log_openLogFile_.tmp"
	_ = os.Remove(testfile)
	var tlog strings.Builder
	now := time.Now().String()
	msg1 := now + " msg1\n"
	msg2 := now + " msg2\n"
	//
	wrc := openLogFile(testfile, &tlog)
	wrc.Write([]byte(msg1))
	wrc.Close()
	//
	wrc = openLogFile(testfile, &tlog)
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
	if tlog.String() != "" {
		t.Error("0xE4CE4D")
	}
	_ = os.Remove(testfile)
}

// go test -run Test_log_openLogFile_2
//
func Test_log_openLogFile_2(t *testing.T) {
	//
	// must fail to open a file with an invalid name
	var tlog strings.Builder
	wrc := openLogFile("\\:/", &tlog)
	if wrc != nil {
		t.Error("0xEC8BA1")
	}
	if !strings.Contains(tlog.String(), "ERROR 0x") {
		t.Error("0xEA1DC3")
	}
}

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
}

// logMakeMessage(tm time.Time, a ...interface{}) string
//
// go test -run Test_log_logMakeMessage_
//
func Test_log_logMakeMessage_(t *testing.T) {
	for _, it := range []struct {
		a    []interface{}
		want string
	}{
		{nil, "@tm "},
		{[]interface{}{"line1\nline2"}, "@tm line1\n@tm line2"},
		{[]interface{}{"abc", 123, "de", 34}, "@tm abc 123 de 34"},
	} {
		var (
			tm   = time.Now()
			tms  = tm.String()[:19]
			got  = logMakeMessage(tm, it.a...)
			want = strings.ReplaceAll(it.want, "@tm", tms)
		)
		if got != want {
			t.Errorf("0xEB14AF"+" logMakeMessage(<%s>, %#v)"+
				"\n want: %#v"+
				"\n  got: %#v",
				tms, it.a, want, got)
		}
	}
}

// end
