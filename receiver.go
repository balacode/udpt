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

	// Config contains UDP and other configuration settings.
	// These settings normally don't need to be changed.
	Config ConfigSettings

	// ReceiveData is a callback function you must specify. This Receiver
	// will call it when a data item has been fully transferred.
	//
	// The 'name' and 'data' parameters will contain the values
	// sent by Sender.Send() or Sender.SendString().
	//
	// The reason there are two parameters is to separate metadata like
	// timestamps or filenames from the bytes of the transferred resource.
	//
	ReceiveData func(name string, data []byte) error

	// ProvideData is a callback function you must specify. This
	// Receiver will call it to read back the named data item.
	//
	// This is needed to send back a confirmation hash to the Sender
	// to confirm the transfer. The Receiver carries out the hashing.
	//
	// If the resource specified by 'name' is found, the function
	// should return its bytes with a nil error value.
	//
	// If the resource is not found, return (nil, nil).
	//
	// If there is some error, return (nil, <error>).
	//
	// NOTE: This callback may be removed because the hashing could be done
	// internally by the Receiver, so there'll be no need for a callback.
	//
	ProvideData func(name string) ([]byte, error)

	// -------------------------------------------------------------------------

	// receivingDataItem contains the data item
	// currently being received from the Sender.
	receivingDataItem dataItem
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
		return ob.logError(0xE1C1A9, ENilReceiver)
	}
	if ob.Port < 1 || ob.Port > 65535 {
		return ob.logError(0xE58B2F, "invalid Receiver.Port:", ob.Port)
	}
	if ob.Config.Cipher == nil {
		var aes aesCipher
		err := aes.InitCipher(ob.CryptoKey)
		if err != nil {
			return ob.logError(0xE36A3C, err)
		}
		ob.Config.Cipher = &aes
	}
	err := ob.Config.Cipher.ValidateKey(ob.CryptoKey)
	if err != nil {
		return ob.logError(0xE3A5FF, "invalid Receiver.CryptoKey:", err)
	}
	err = ob.Config.Validate()
	if err != nil {
		return ob.logError(0xE14BC8, err)
	}
	if ob.ReceiveData == nil {
		return ob.logError(0xE82C9E, "nil Receiver.ReceiveData")
	}
	if ob.ProvideData == nil {
		return ob.logError(0xE4E2C1, "nil Receiver.ProvideData")
	}
	if ob.Config.VerboseReceiver {
		ob.logInfo(strings.Repeat("-", 80))
		ob.logInfo("UDPT started in receiver mode")
	}
	udpAddr, err := net.ResolveUDPAddr("udp",
		fmt.Sprintf("0.0.0.0:%d", ob.Port))
	if err != nil {
		return ob.logError(0xE1D68C, err)
	}
	conn, err := net.ListenUDP("udp", udpAddr) // (*net.UDPConn, error)
	if err != nil {
		return ob.logError(0xEBF95F, err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			_ = ob.logError(0xE15F3A, err)
		}
	}()
	if ob.Config.VerboseReceiver {
		ob.logInfo("Receiver.Run() called net.ListenUDP")
	}
	encryptedReq := make([]byte, ob.Config.PacketSizeLimit)
	for conn != nil {
		// 'encryptedReq' is overwritten after every readFromUDPConn
		nRead, addr, err :=
			readFromUDPConn(conn, encryptedReq, ob.Config.ReplyTimeout)
		if err != nil {
			_ = ob.logError(0xEA288A, err)
			continue
		}
		recv, err := ob.Config.Cipher.Decrypt(encryptedReq[:nRead])
		if err != nil {
			_ = ob.logError(0xE7D2C4, err)
			continue
		}
		if ob.Config.VerboseReceiver {
			ob.logInfo()
			ob.logInfo(strings.Repeat("-", 80))
			ob.logInfo("Receiver read", nRead, "bytes from", addr)
		}
		var reply []byte
		switch {
		case len(recv) == 0:
			_ = ob.logError(0xE6B3BA, "received no data")
			continue
		case bytes.HasPrefix(recv, []byte(tagDataItemHash)):
			reply, err = ob.sendDataItemHash(recv)
			if err != nil {
				_ = ob.logError(0xE98D72, err)
				continue
			}
		case bytes.HasPrefix(recv, []byte(tagFragment)):
			reply, err = ob.receiveFragment(recv)
			if err != nil {
				_ = ob.logError(0xE3A46C, err)
				continue
			}
		default:
			_ = ob.logError(0xE985CC, "invalid packet header")
			reply = []byte("invalid_packet_header")
		}
		encryptedReply, err := ob.Config.Cipher.Encrypt(reply)
		if err != nil {
			_ = ob.logError(0xE6E8C7, err)
			continue
		}
		deadline := time.Now().Add(ob.Config.WriteTimeout)
		err = conn.SetWriteDeadline(deadline)
		if err != nil {
			_ = ob.logError(0xE1F2C4, err)
			continue
		}
		nWrit, err := conn.WriteTo(encryptedReply, addr)
		if err != nil {
			_ = ob.logError(0xEA63C4, err)
			continue
		}
		if ob.Config.VerboseReceiver {
			ob.logInfo("Receiver wrote", nWrit, "bytes to", addr)
		}
	}
	return nil
} //                                                                         Run

// -----------------------------------------------------------------------------
// # Packet Handlers

// receiveFragment handles a tagFragment packet sent by a Sender, and
// sends back a confirmation packet (tagConfirmation) to the Sender.
func (ob *Receiver) receiveFragment(recv []byte) ([]byte, error) {
	if ob == nil {
		return nil, ob.logError(0xE6CD62, ENilReceiver)
	}
	if !bytes.HasPrefix(recv, []byte(tagFragment)) {
		return nil, ob.logError(0xE4F3C5, "missing header")
	}
	err := ob.Config.Validate()
	if err != nil {
		return nil, ob.logError(0xE9B5C7, err)
	}
	dataOffset := bytes.Index(recv, []byte("\n"))
	if dataOffset == -1 {
		return nil, ob.logError(0xE6CF52, "newline not found")
	}
	dataOffset++ // skip newline
	var (
		header  = string(recv[len(tagFragment):dataOffset])
		name    = getPart(header, "name:", " ")
		hexHash = getPart(header, "hash:", " ")
		sn      = getPart(header, "sn:", " ")
		count   = getPart(header, "count:", "\n")
	)
	index, err := strconv.Atoi(sn)
	if err != nil {
		return nil, ob.logError(0xE14D6A, "bad 'sn':", sn)
	}
	index--
	packetCount, err := strconv.Atoi(count)
	if err != nil {
		return nil, ob.logError(0xE76D48, "bad 'count'")
	}
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return nil, ob.logError(0xE5CA62, err)
	}
	if packetCount < 1 {
		return nil, ob.logError(0xE29C5B, "invalid packetCount:", packetCount)
	}
	if index < 0 || index >= packetCount {
		return nil, ob.logError(0xE8CF4D,
			"index", index, "out of range 0 -", packetCount-1)
	}
	it := &ob.receivingDataItem
	it.Retain(name, hash, packetCount)
	compressedData := recv[dataOffset:]
	if len(compressedData) < 1 {
		return nil, ob.logError(0xE92B0F, "received no data")
	}
	// store the current piece
	if len(it.CompressedPieces[index]) == 0 {
		it.CompressedPieces[index] = compressedData
	} else if !bytes.Equal(compressedData, it.CompressedPieces[index]) {
		return nil, ob.logError(0xE981DA, "unknown packet alteration")
	}
	if it.IsLoaded() {
		if ob.ReceiveData == nil {
			return nil, ob.logError(0xE49E2A, "nil Receiver.ReceiveData")
		}
		data, err := it.UnpackBytes()
		if err != nil {
			return nil, ob.logError(0xE3DB1D, err)
		}
		err = ob.ReceiveData(it.Name, data)
		if err != nil {
			return nil, ob.logError(0xE9BD1B, err)
		}
		ob.logInfo("received:", it.Name)
		if ob.Config.VerboseReceiver {
			it.PrintInfo("receiveFragment", ob.logInfo)
		}
		it.Reset()
	}
	confirmedHash, err := getHash(recv)
	if err != nil {
		return nil, ob.logError(0xE0B57C, err)
	}
	reply := append([]byte(tagConfirmation), confirmedHash...)
	return reply, nil
} //                                                             receiveFragment

// sendDataItemHash handles a tagDataItemHash sent by a Sender.
func (ob *Receiver) sendDataItemHash(req []byte) ([]byte, error) {
	if ob == nil {
		return nil, ob.logError(0xE24A7B, ENilReceiver)
	}
	if !bytes.HasPrefix(req, []byte(tagDataItemHash)) {
		return nil, ob.logError(0xE7B653, "missing header")
	}
	if ob.ProvideData == nil {
		return nil, ob.logError(0xE73A1C, "nil ProvideData")
	}
	name := string(req[len(tagDataItemHash):])
	data, err := ob.ProvideData(name)
	if err != nil {
		return nil, ob.logError(0xE7F7C9, err)
	}
	hash, err := getHash(data)
	if err != nil {
		return nil, ob.logError(0xE2F3D5, err)
	}
	reply := []byte(tagDataItemHash + fmt.Sprintf("%X", hash))
	return reply, nil
} //                                                            sendDataItemHash

// -----------------------------------------------------------------------------
// # Logging Methods

// logError returns a new error value generated by joining id and args
// and optionally calls Receiver.LogFunc (if not nil) to log the error.
func (ob *Receiver) logError(id uint32, args ...interface{}) error {
	ret := makeError(id, args...)
	if ob.Config.LogFunc != nil {
		msg := ret.Error()
		ob.Config.LogFunc(msg)
	}
	return ret
} //                                                                    logError

// logInfo calls Receiver.LogFunc (if not nil) to log a message.
func (ob *Receiver) logInfo(args ...interface{}) {
	if ob.Config.LogFunc != nil {
		ob.Config.LogFunc(args...)
	}
} //                                                                     logInfo

// end
