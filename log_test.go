// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[log_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run _log_

// -----------------------------------------------------------------------------

// LogPrint()
//
// go test -run _log_LogPrint_
//
func Test_log_LogPrint_(t *testing.T) {
	LogPrint()
} //                                                          Test_log_LogPrint_

// MakeLogFunc(printMsg bool, logFile string) func(args ...interface{})
//
// go test -run _log_MakeLogFunc_
//
func Test_log_MakeLogFunc_(t *testing.T) {
	//
	// prepare: delete log file and mock time.Now()
	logFile := "udpt.Test_log_MakeLogFunc_.tmp.log"
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
} //                                                       Test_log_MakeLogFunc_

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
