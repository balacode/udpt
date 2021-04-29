// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[receiver.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Receiver _ _
type Receiver struct {
	currentDataItem DataItem
	writeDataFn     func(name string, data []byte) error
	readDataFn      func(name string) ([]byte, error)
} //                                                                    Receiver

// RunReceiver starts a goroutine that runs the receiver in an infinite loop.
func RunReceiver(
	writeDataFn func(name string, data []byte) error,
	readDataFn func(name string) ([]byte, error),
) {
	go func() {
		receiver := Receiver{
			writeDataFn: writeDataFn,
			readDataFn:  readDataFn,
		}
		err := receiver.run()
		if err != nil {
			_ = logError(0xE0A4AC, err)
		}
	}()
} //                                                                 RunReceiver

// -----------------------------------------------------------------------------
// # Main Loop

// run _ _
func (ob *Receiver) run() error {
	if ob == nil {
		return logError(0xE1C1A9, ":", ENilReceiver)
	}
	err := Config.Validate()
	if err != nil {
		return logError(0xE14BC8, err)
	}
	if Config.VerboseReceiver {
		logInfo(strings.Repeat("-", 80))
		logInfo("UDPT started in receiver mode")
	}
	udpAddr, err := net.ResolveUDPAddr("udp",
		fmt.Sprintf("0.0.0.0:%d", Config.Port))
	if err != nil {
		return logError(0xE1D68C, "(ResolveUDPAddr):", err)
	}
	conn, err := net.ListenUDP("udp", udpAddr) // (*net.UDPConn, error)
	if err != nil {
		return logError(0xEBF95F, "(ListenPacket):", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			_ = logError(0xE15F3A, "(Close):", err)
		}
	}()
	if Config.VerboseReceiver {
		logInfo("Receiver.run() called net.ListenUDP")
	}
	encryptedReq := make([]byte, Config.PacketSizeLimit)
	for {
		// 'encryptedReq' is overwritten after every readFromUDPConn
		nRead, addr, err := readFromUDPConn(conn, encryptedReq)
		if err != nil {
			_ = logError(0xEA288A, "(readFromUDPConn):", err)
			continue
		}
		recv, err := aesDecrypt(encryptedReq[:nRead], Config.AESKey)
		if err != nil {
			_ = logError(0xE7D2C4, "(aesDecrypt):", err)
			continue
		}
		if Config.VerboseReceiver {
			logInfo()
			logInfo(strings.Repeat("-", 80))
			logInfo("Receiver read", nRead, "bytes from", addr)
		}
		var reply []byte
		switch {
		case len(recv) == 0:
			_ = logError(0xE6B3BA, ": received no data")
			continue
		case bytes.HasPrefix(recv, []byte(DATA_ITEM_HASH)):
			reply, err = ob.sendDataItemHash(recv)
			if err != nil {
				_ = logError(0xE98D72, "(sendDataItemHash):", err)
				continue
			}
		case bytes.HasPrefix(recv, []byte(FRAGMENT)):
			reply, err = ob.receiveFragment(recv)
			if err != nil {
				_ = logError(0xE3A46C, "(receiveFragment):", err)
				continue
			}
		default:
			_ = logError(0xE985CC, ": Invalid packet header")
			reply = []byte("invalid_packet_header")
		}
		encryptedReply, err := aesEncrypt(reply, Config.AESKey)
		if err != nil {
			_ = logError(0xE6E8C7, "(aesEncrypt):", err)
			continue
		}
		deadline := time.Now().Add(Config.WriteTimeout)
		err = conn.SetWriteDeadline(deadline)
		if err != nil {
			_ = logError(0xE1F2C4, "(SetWriteDeadline):", err)
			continue
		}
		nWrit, err := conn.WriteTo(encryptedReply, addr)
		if err != nil {
			_ = logError(0xEA63C4, "(WriteTo):", err)
			continue
		}
		if Config.VerboseReceiver {
			logInfo("Receiver wrote", nWrit, "bytes to", addr)
		}
	}
} //                                                                         run

// -----------------------------------------------------------------------------
// # Handler Methods

// receiveFragment _ _
func (ob *Receiver) receiveFragment(recv []byte) ([]byte, error) {
	if ob == nil {
		return nil, logError(0xE6CD62, ":", ENilReceiver)
	}
	if !bytes.HasPrefix(recv, []byte(FRAGMENT)) {
		return nil, logError(0xE4F3C5, ": missing header")
	}
	err := Config.Validate()
	if err != nil {
		return nil, logError(0xE9B5C7, err)
	}
	dataOffset := bytes.Index(recv, []byte("\n"))
	if dataOffset == -1 {
		return nil, logError(0xE6CF52, ": newline not found")
	}
	dataOffset++ // skip newline
	var (
		header  = string(recv[len(FRAGMENT):dataOffset])
		name    = getPart(header, "name:", " ")
		hexHash = getPart(header, "hash:", " ")
		sn      = getPart(header, "sn:", " ")
		count   = getPart(header, "count:", "\n")
	)
	index, err := strconv.Atoi(sn)
	if err != nil {
		return nil, logError(0xE14D6A, ": bad 'sn':", sn)
	}
	index--
	packetCount, err := strconv.Atoi(count)
	if err != nil {
		return nil, logError(0xE76D48, ": bad 'count'")
	}
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return nil, logError(0xE5CA62, "(hex.DecodeString):", err)
	}
	if packetCount < 1 {
		return nil, logError(0xE29C5B, ": Invalid packetCount:", packetCount)
	}
	if index < 0 || index >= packetCount {
		return nil, logError(0xE8CF4D, ":",
			"index", index, "out of range 0 -", packetCount-1)
	}
	it := &ob.currentDataItem
	it.Retain(name, hash, packetCount)
	compressedData := recv[dataOffset:]
	if len(compressedData) < 1 {
		return nil, logError(0xE92B0F, ": received no data")
	}
	// store the current piece
	if len(it.CompressedPieces[index]) == 0 {
		it.CompressedPieces[index] = compressedData
	} else if !bytes.Equal(compressedData, it.CompressedPieces[index]) {
		return nil, logError(0xE981DA, ": unknown packet change")
	}
	if it.IsLoaded() {
		if ob.writeDataFn == nil {
			return nil, logError(0xE49E2A, "writeDataFn is nil")
		}
		data, err := it.UnpackBytes()
		if err != nil {
			return nil, logError(0xE3DB1D, "(UnpackBytes):", err)
		}
		err = ob.writeDataFn(it.Name, data)
		if err != nil {
			return nil, logError(0xE9BD1B, "(writeDataFn):", err)
		}
		logInfo("received:", it.Name)
		if Config.VerboseReceiver {
			it.PrintInfo("receiveFragment")
		}
		it.Reset()
	}
	confirmedHash := getHash(recv)
	reply := append([]byte(FRAGMENT_CONFIRMATION), confirmedHash...)
	return reply, nil
} //                                                             receiveFragment

// sendDataItemHash handles a DATA_ITEM_HASH sent by a Sender.
func (ob *Receiver) sendDataItemHash(req []byte) ([]byte, error) {
	if ob == nil {
		return nil, logError(0xE24A7B, ":", ENilReceiver)
	}
	if !bytes.HasPrefix(req, []byte(DATA_ITEM_HASH)) {
		return nil, logError(0xE7B653, ": missing header")
	}
	if ob.readDataFn == nil {
		return nil, logError(0xE73A1C, "readDataFn is nil")
	}
	name := string(req[len(DATA_ITEM_HASH):])
	data, err := ob.readDataFn(name)
	if err != nil {
		return nil, logError(0xE7F7C9, "(readDataFn):", err)
	}
	hash := getHash(data)
	reply := []byte(DATA_ITEM_HASH + fmt.Sprintf("%X", hash))
	return reply, nil
} //                                                            sendDataItemHash

// end
