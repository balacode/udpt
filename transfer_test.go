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

// go test --run Transfer1
func Test_Transfer1_(t *testing.T) {
	fmt.Println("Test_Transfer1_")
	const itemCount = 50
	const itemSize = 100
	testTransfer(itemCount, itemSize, t)
} //                                                             Test_Transfer1_

// go test --run Transfer2
func Test_Transfer2_(t *testing.T) {
	fmt.Println("Test_Transfer2_")
	const itemCount = 10
	const itemSize = 10 * 1024 * 1024 // 10 MiB
	testTransfer(itemCount, itemSize, t)
} //                                                             Test_Transfer2_

// go test --run Transfer3
func Test_Transfer3_(t *testing.T) {
	fmt.Println("Test_Transfer3_")
	const itemCount = 1
	const itemSize = 100 * 1024 * 1024 // 100 MiB
	testTransfer(itemCount, itemSize, t)
} //                                                             Test_Transfer3_

// testTransfer runs a transfer test with different packet counts and sizes.
//
// This test sends several packets from a Sender to a Receiver.
// After sending, it checks if all the packets have been delivered.
//
// itemCount specifies the number of b
//
//
func testTransfer(itemCount, itemSize int, t *testing.T) {
	var N = itemCount
	//
	var cryptoKey = []byte("aA2Xh41FiC4Wtj3e5b2LbytMdn6on7P0")
	//
	// this map collects received keys and values
	received := make(map[string][]byte, N)
	//
	// disable log buffering and enable verbose logging: for demos/debugging
	cfg := NewDefaultConfig()
	if false {
		cfg.LogFunc = LogPrint
		cfg.VerboseSender = true
		cfg.VerboseReceiver = true
	}
	// set-up and run the receiver
	receiver := Receiver{
		Port: 1234, CryptoKey: cryptoKey, Config: cfg,
		//
		ReceiveData: func(name string, data []byte) error {
			k, v := name, data
			received[k] = []byte(v)
			return nil
		},
		ProvideData: func(name string) ([]byte, error) {
			v := received[name]
			return v, nil
		},
	}
	go func() { _ = receiver.Run() }()
	defer func() { receiver.Stop() }()
	//
	// make a map of N messages
	for i := 0; i < N; i++ {
		k := fmt.Sprint("P", i)
		v := fmt.Sprintf("%04d", i)
		received[k] = []byte(v)
	}
	// send the messages to the receiver
	time.Sleep(time.Second)
	sender := Sender{
		Address: "127.0.0.1", Port: 1234, CryptoKey: cryptoKey, Config: cfg,
	}
	makeKV := func(i int) (string, []byte) {
		sn := fmt.Sprintf("%04d", i)
		k := fmt.Sprint("msg", i)
		v := strings.Repeat(sn, (itemSize/4)+1)[:itemSize]
		return k, []byte(v)
	}
	for i := 0; i < N; i++ {
		k, v := makeKV(i)
		err := sender.Send(k, []byte(v))
		if err != nil {
			t.Error("failed sending "+k+":", err)
		}
	}
	time.Sleep(time.Second)
	//
	// compare received to expected values
	for i := 0; i < N; i++ {
		k, vS := makeKV(i)
		vR := received[k]
		if !bytes.Equal(vS, vR) {
			t.Error("mismatch for key:", k,
				"len(vS):", len(vS),
				"len(vR):", len(vR))
		}
	}
	if false {
		cfg.LogFunc = LogPrint
		cfg.VerboseSender = true
		cfg.VerboseReceiver = true
		sender.PrintInfo()
	}
	time.Sleep(time.Second)
} //                                                                testTransfer

// end
