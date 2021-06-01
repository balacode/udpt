// -----------------------------------------------------------------------------
// github.com/balacode/udpt                             /demo/[transfer_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run Test_transfer_*

// -----------------------------------------------------------------------------

// go test -run Test_transfer_1
//
func Test_transfer_1(t *testing.T) {
	const itemCount = 50
	const itemSize = 100
	testTransfer(itemCount, itemSize, t)
}

// go test -run Test_transfer_2
//
func Test_transfer_2(t *testing.T) {
	const itemCount = 10
	const itemSize = 10 * 1024 * 1024 // 10 MiB
	testTransfer(itemCount, itemSize, t)
}

// go test -run Test_transfer_3
//
func Test_transfer_3(t *testing.T) {
	const itemCount = 1
	const itemSize = 100 * 1024 * 1024 // 100 MiB
	testTransfer(itemCount, itemSize, t)
}

// testTransfer runs a transfer test with different packet counts and sizes.
//
// This test sends several packets from a Sender to a Receiver.
// After sending, it checks if all the packets have been delivered.
//
// itemCount specifies the number of data items to send
//
// itemSize specifies the size of each message in bytes
//
func testTransfer(itemCount, itemSize int, t *testing.T) {
	var N = itemCount
	//
	var cryptoKey = []byte("aA2Xh41FiC4Wtj3e5b2LbytMdn6on7P0")
	//
	// this map collects received keys and values
	received := make(map[string][]byte, N)
	//
	// enable verbose logging but don't print the output
	cf := NewDefaultConfig()
	cf.LogWriter = nil
	cf.VerboseSender = true
	cf.VerboseReceiver = true
	//
	// set-up and run the receiver
	rc := Receiver{
		Port: 9876, CryptoKey: cryptoKey, Config: cf,
		//
		ReceiveData: func(k string, v []byte) error {
			received[k] = []byte(v)
			return nil
		},
		ProvideData: func(k string) ([]byte, error) {
			v := received[k]
			return v, nil
		},
	}
	go func() { _ = rc.Run() }()
	defer func() { rc.Stop() }()
	//
	// make a map of N messages
	for i := 0; i < N; i++ {
		k := fmt.Sprint("P", i)
		v := fmt.Sprintf("%04d", i)
		received[k] = []byte(v)
	}
	// send the messages to the receiver
	time.Sleep(time.Second)
	sd := Sender{
		Address: "127.0.0.1:9876", CryptoKey: cryptoKey, Config: cf,
	}
	makeKV := func(i int) (string, []byte) {
		sn := fmt.Sprintf("%04d", i)
		k := fmt.Sprint("msg", i)
		v := strings.Repeat(sn, (itemSize/4)+1)[:itemSize]
		return k, []byte(v)
	}
	for i := 0; i < N; i++ {
		k, v := makeKV(i)
		err := sd.Send(k, []byte(v))
		if err != nil {
			t.Error("0xEA8C61", "failed sending "+k+":", err)
		}
	}
	time.Sleep(time.Second)
	//
	// compare received to expected values
	for i := 0; i < N; i++ {
		k, vS := makeKV(i)
		vR := received[k]
		if !bytes.Equal(vS, vR) {
			t.Error("0xE67F41", "mismatch for key:", k,
				"len(vS):", len(vS),
				"len(vR):", len(vR))
		}
	}
	time.Sleep(time.Second)
}

// end
