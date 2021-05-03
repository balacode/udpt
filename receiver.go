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

// Receiver receives data items sent by Send() or SendString().
type Receiver struct {

	// Port is the port number of the listening server.
	// This number must be between 1 and 65535.
	Port int

	// CryptoKey is the secret symmetric encryption key that
	// must be shared between the sender and the receiver.
	// The correct size of this key depends
	// on the implementation of SymmetricCipher.
	CryptoKey []byte

	// Cipher is the object that handles enryption and decryption.
	//
	// It must implement the SymmetricCipher interface which is defined in
	// this package. If you don't specify Cipher, then encyption will
	// be done using AESCipher, the default encryption used in this package.
	//
	Cipher SymmetricCipher

	// Config _ _
	Config ConfigSettings

	// ReceiveData _ _
	ReceiveData func(name string, data []byte) error

	// ProvideData _ _
	ProvideData func(name string) ([]byte, error)

	// receivingDataItem _ _
	receivingDataItem dataItemStruct
} //                                                                    Receiver

// -----------------------------------------------------------------------------
// # Public Methods

// Run runs the receiver in a loop to process incoming packets.
//
// It calls ReceiveData when a data transfer is complete, after the
// receiver has received, decrypted and re-assembled a data item.
//
// It calls ProvideData when it needs to calculate the hash of a
// previously-received data item. This hash is sent to the sender
// so it can confirm that a data transfer is successful.
//
func (ob *Receiver) Run() error {
	if ob == nil {
		return logError(0xE1C1A9, ":", ENilReceiver)
	}
	if ob.Port < 1 || ob.Port > 65535 {
		return logError(0xE58B2F, "invalid Port:", ob.Port)
	}
	if ob.Cipher == nil {
		var aes AESCipher
		err := aes.InitCipher(ob.CryptoKey)
		if err != nil {
			return logError(0xE36A3C, "(aes.InitCipher):", err)
		}
		ob.Cipher = &aes
	}
	err := ob.Cipher.ValidateKey(ob.CryptoKey)
	if err != nil {
		return logError(0xE3A5FF, "invalid Receiver.CryptoKey:", err)
	}
	err = ob.Config.Validate()
	if err != nil {
		return logError(0xE14BC8, err)
	}
	if ob.ReceiveData == nil {
		return logError(0xE82C9E, ": ReceiveData func. is nil.")
	}
	if ob.ProvideData == nil {
		return logError(0xE4E2C1, ": ProvideData func. is nil.")
	}
	if ob.Config.VerboseReceiver {
		logInfo(strings.Repeat("-", 80))
		logInfo("UDPT started in receiver mode")
	}
	udpAddr, err := net.ResolveUDPAddr("udp",
		fmt.Sprintf("0.0.0.0:%d", ob.Port))
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
	if ob.Config.VerboseReceiver {
		logInfo("Receiver.Run() called net.ListenUDP")
	}
	encryptedReq := make([]byte, ob.Config.PacketSizeLimit)
	for conn != nil {
		// 'encryptedReq' is overwritten after every readFromUDPConn
		nRead, addr, err :=
			readFromUDPConn(conn, encryptedReq, ob.Config.ReplyTimeout)
		if err != nil {
			_ = logError(0xEA288A, "(readFromUDPConn):", err)
			continue
		}
		recv, err := ob.Cipher.Decrypt(encryptedReq[:nRead])
		if err != nil {
			_ = logError(0xE7D2C4, "(Decrypt):", err)
			continue
		}
		if ob.Config.VerboseReceiver {
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
		encryptedReply, err := ob.Cipher.Encrypt(reply)
		if err != nil {
			_ = logError(0xE6E8C7, "(Encrypt):", err)
			continue
		}
		deadline := time.Now().Add(ob.Config.WriteTimeout)
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
		if ob.Config.VerboseReceiver {
			logInfo("Receiver wrote", nWrit, "bytes to", addr)
		}
	}
	return nil
} //                                                                         Run

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
	err := ob.Config.Validate()
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
	it := &ob.receivingDataItem
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
		if ob.ReceiveData == nil {
			return nil, logError(0xE49E2A, "ReceiveData is nil")
		}
		data, err := it.UnpackBytes()
		if err != nil {
			return nil, logError(0xE3DB1D, "(UnpackBytes):", err)
		}
		err = ob.ReceiveData(it.Name, data)
		if err != nil {
			return nil, logError(0xE9BD1B, "(ReceiveData):", err)
		}
		logInfo("received:", it.Name)
		if ob.Config.VerboseReceiver {
			it.PrintInfo("receiveFragment", logInfo)
		}
		it.Reset()
	}
	confirmedHash, err := getHash(recv)
	if err != nil {
		return nil, logError(0xE0B57C, "(getHash):", err)
	}
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
	if ob.ProvideData == nil {
		return nil, logError(0xE73A1C, "ProvideData is nil")
	}
	name := string(req[len(DATA_ITEM_HASH):])
	data, err := ob.ProvideData(name)
	if err != nil {
		return nil, logError(0xE7F7C9, "(ProvideData):", err)
	}
	hash, err := getHash(data)
	if err != nil {
		return nil, logError(0xE2F3D5, "(getHash):", err)
	}
	reply := []byte(DATA_ITEM_HASH + fmt.Sprintf("%X", hash))
	return reply, nil
} //                                                            sendDataItemHash

// end
