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

// -----------------------------------------------------------------------------

// to run all tests in this file:
// go test -v -run _log_

// pl(args ...interface{})
//
// go test -run _log_pl_
//
func Test_log_pl_(t *testing.T) {
	pl()
} //                                                                Test_log_pl_

// LogPrint(args ...interface{})
//
// go test -run _log_LogPrint_
//
func Test_log_LogPrint_(t *testing.T) {
	LogPrint()
} //                                                          Test_log_LogPrint_

// MakeLogFunc(printMsg bool, logFile string) func(args ...interface{})
//
// go test -run _log_MakeLogFunc_1_
//
func Test_log_MakeLogFunc_1_(t *testing.T) {
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
		t.Error("MakeLogFunc returned nil")
	}
	fn("abc", 123, "45", "de")
	//
	// check results
	time.Sleep(500 * time.Millisecond) // wait for logger finish writing
	data, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Error("Failed to read log file")
		return
	}
	got := string(data)
	expect := "@tm abc 123 45 de\n"
	expect = strings.ReplaceAll(expect, "@tm", tm.String()[:19])
	//
	if got != expect {
		t.Errorf("Wrong text in log file "+
			"\n expect: %#v"+
			"\n    got: %#v",
			expect, got)
	}
	// cleanup
	_ = os.Remove(logFile)
	logTimeNow = time.Now
} //                                                     Test_log_MakeLogFunc_1_

// MakeLogFunc(printMsg bool, logFile string) func(args ...interface{})
//
// go test -run _log_MakeLogFunc_2_
//
func Test_log_MakeLogFunc_2_(t *testing.T) {
	const filename = "udpt.Test_log_MakeLogFunc_.tmp"
	const msg = "test message #5067124389"
	_ = os.Remove(filename)
	fn := MakeLogFunc(false, filename)
	fn(msg)
	time.Sleep(time.Second)
	data, err := ioutil.ReadFile(filename)
	if !strings.Contains(string(data), msg) || err != nil {
		t.Error("0xEC7ED0")
	}
	_ = os.Remove(filename)
} //                                                     Test_log_MakeLogFunc_2_

// -----------------------------------------------------------------------------

// (ob *logEntry) outputDI(
//     con io.Writer,
//     openLogFile func(filename string, con io.Writer) io.WriteCloser,
// )
//
// go test -run _log_logEntry_outputDI_
//
func Test_log_logEntry_outputDI_(t *testing.T) {
	{
		// test writing to console without logging to a file
		const msg = "test message #3146250897"
		l := logEntry{printMsg: true, logFile: "", msg: msg}
		var sb strings.Builder
		l.outputDI(&sb, nil)
		if sb.String() != msg+"\n" {
			t.Error("0xE58FB6")
		}
	}
	const filename = "udpt.Test_log_logEntry_outputDI_.tmp"
	_ = os.Remove(filename)
	{
		// test logging to a file without writing to console
		const msg = "test message #9473258061"
		l := logEntry{printMsg: false, logFile: filename, msg: msg}
		var sb strings.Builder
		l.outputDI(&sb, openLogFile)
		if sb.String() != "" {
			t.Error("0xE8BB8F")
		}
		data, err := ioutil.ReadFile(filename)
		if string(data) != msg+"\n" || err != nil {
			t.Error("0xED98CA")
		}
	}
	_ = os.Remove(filename)
	{
		// test that no file is created when openLogFile() returns nil
		const msg = "test message #2361498057"
		l := logEntry{printMsg: false, logFile: filename, msg: msg}
		openLogFile := func(string, io.Writer) io.WriteCloser { return nil }
		var sb strings.Builder
		l.outputDI(&sb, openLogFile)
		data, err := ioutil.ReadFile(filename)
		if data != nil || !os.IsNotExist(err) {
			t.Error("0xE7BC99")
		}
	}
	_ = os.Remove(filename)
	{
		// test that errors are written to console when Write() and Close() fail
		const msg = "test message #1540938672"
		l := logEntry{printMsg: false, logFile: filename, msg: msg}
		openLogFile := func(string, io.Writer) io.WriteCloser {
			return &mockWriteCloser{failWrite: true, failClose: true}
		}
		var sb strings.Builder
		l.outputDI(&sb, openLogFile)
		//
		// console must contain two error messages
		con := sb.String()
		if !strings.Contains(con, "ERROR 0x") ||
			!strings.Contains(con, "failed mockWriteCloser.Write") ||
			!strings.Contains(con, "failed mockWriteCloser.Close") {
			t.Error("0xEA1B28")
		}
		// must not create/write to file
		data, err := ioutil.ReadFile(filename)
		if data != nil || !os.IsNotExist(err) {
			t.Error("0xEB38C2")
		}
	}
	_ = os.Remove(filename)
} //                                                 Test_log_logEntry_outputDI_

// openLogFile(filename string, con io.Writer) io.WriteCloser
//
// go test -run _log_openLogFile_
//
func Test_log_openLogFile_(t *testing.T) {
	{
		// must succeed opening a file and writing to it two times
		const filename = "udpt.Test_log_openLogFile_.tmp"
		_ = os.Remove(filename)
		var sb strings.Builder
		now := time.Now().String()
		msg1 := now + " msg1\n"
		msg2 := now + " msg2\n"
		//
		wrc := openLogFile(filename, &sb)
		wrc.Write([]byte(msg1))
		wrc.Close()
		//
		wrc = openLogFile(filename, &sb)
		wrc.Write([]byte(msg2))
		wrc.Close()
		//
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Error("0xE7B84B")
		}
		content := string(data)
		if content != msg1+msg2 {
			t.Error("0xE24C06")
		}
		if sb.String() != "" {
			t.Error("0xE4CE4D")
		}
		_ = os.Remove(filename)
	}
	{
		// must fail to open a file with an invalid name
		var sb strings.Builder
		wrc := openLogFile("\\:/", &sb)
		if wrc != nil {
			t.Error("0xEC8BA1")
		}
		if !strings.Contains(sb.String(), "ERROR 0x") {
			t.Error("0xEA1DC3")
		}
	}
} //                                                       Test_log_openLogFile_

// logInit()
//
// go test -run _log_logInit_
//
func Test_log_logInit_(t *testing.T) {
	logInit()
	if logChan == nil {
		t.Error("logChan not initialized")
	}
} //                                                           Test_log_logInit_

// logMakeMessage(args ...interface{}) string
//
// go test -run _log_logMakeMessage_
//
func Test_log_logMakeMessage_(t *testing.T) {
	for _, it := range []struct {
		args   []interface{}
		expect string
	}{
		{nil, "@tm "},
		{[]interface{}{"line1\nline2"}, "@tm line1\n@tm line2"},
		{[]interface{}{"abc", 123, "de", 34}, "@tm abc 123 de 34"},
	} {
		var (
			tm     = time.Now()
			tms    = tm.String()[:19]
			got    = logMakeMessage(tm, it.args...)
			expect = strings.ReplaceAll(it.expect, "@tm", tms)
		)
		if got != expect {
			t.Errorf("logMakeMessage(<%s>, %#v)"+
				"\n expect: %#v"+
				"\n    got: %#v",
				tms, it.args, expect, got)
		}
	}
} //                                                    Test_log_logMakeMessage_

// end
