// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[sender_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"strings"
	"testing"
)

// to run all tests in this file:
// go test -v -run Test_Sender_

// -----------------------------------------------------------------------------

// (ob *Sender) LogStats(logFunc ...interface{})
//
// go test -run Test_Sender_LogStats_
//
func Test_Sender_LogStats_(t *testing.T) {
	var sb strings.Builder
	fmtPrintln := func(v ...interface{}) (int, error) {
		sb.WriteString(fmt.Sprintln(v...))
		return 0, nil
	}
	logPrintln := func(v ...interface{}) {
		sb.WriteString(fmt.Sprintln(v...))
	}
	test := func(logFunc interface{}) {
		sb.Reset()
		//
		var snd Sender
		snd.Config = NewDefaultConfig()
		snd.Config.LogFunc = logPrintln
		snd.packets = []senderPacket{{
			sentHash:      []byte{0},
			confirmedHash: []byte{0},
		}}
		snd.LogStats(logFunc)
		//
		got := sb.String()
		expect := "" +
			"SN: 0    T0: 0001-01-01 00:00:00 +000 T1: NONE ✔ 0.0 ms\n" +
			"B. delivered: 0\n" +
			"Bytes lost  : 0\n" +
			"P. delivered: 0\n" +
			"Packets lost: 0\n" +
			"Time in item: 0.0 s\n" +
			"Avg./ Packet: 0.0 ms\n" +
			"Trans. speed: 0.0 KiB/s\n"
		//
		if got != expect {
			t.Error("\n" + "expect:\n" + expect + "\n" + "got:\n" + got)
			fmt.Println([]byte(expect))
			fmt.Println([]byte(got))
		}
	}
	test(nil)
	test(logPrintln)
	test(fmtPrintln)
} //                                                       Test_Sender_LogStats_

// end
