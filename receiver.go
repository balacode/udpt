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
	Config *Configuration

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

	// conn is the UDP connection on which Receiver listens;
	// setting this to nil allows Run() to stop listening
	conn *net.UDPConn

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
	if ob.Config == nil {
		ob.Config = NewDefaultConfig()
	}
	// setup cipher
	if ob.Config.Cipher == nil {
		return ob.logError(0xE62F4F, "nil Receiver.Config.Cipher")
	}
	err := ob.Config.Cipher.SetKey(ob.CryptoKey)
	if err != nil {
		return ob.logError(0xE8A5C6, "invalid Receiver.CryptoKey:", err)
	}
	// check settings
	err = ob.Config.Validate()
	if err != nil {
		return ob.logError(0xE14BC8, err)
	}
	if ob.Port < 1 || ob.Port > 65535 {
		return ob.logError(0xE58B2F, "invalid Receiver.Port:", ob.Port)
	}
	if ob.ReceiveData == nil {
		return ob.logError(0xE82C9E, "nil Receiver.ReceiveData")
	}
	if ob.ProvideData == nil {
		return ob.logError(0xE48CC6, "nil Receiver.ProvideData")
	}
	// prepare for reception
	if ob.Config.VerboseReceiver {
		ob.logInfo(strings.Repeat("-", 80))
		ob.logInfo("UDPT started in receiver mode")
	}
	udpAddr, err := net.ResolveUDPAddr("udp",
		fmt.Sprintf("0.0.0.0:%d", ob.Port))
	if err != nil {
		return ob.logError(0xE1D68C, err)
	}
	ob.conn, err = net.ListenUDP("udp", udpAddr) // (*net.UDPConn, error)
	if err != nil {
		return ob.logError(0xEBF95F, err)
	}
	defer func() {
		if ob.conn != nil {
			err := ob.conn.Close()
			if err != nil {
				_ = ob.logError(0xEC82DB, err)
			}
		}
	}()
	if ob.Config.VerboseReceiver {
		ob.logInfo("Receiver.Run() called net.ListenUDP")
	}
	// receive transmissions
	encryptedReq := make([]byte, ob.Config.PacketSizeLimit)
	for ob.conn != nil {
		// 'encryptedReq' is overwritten after every readFromUDPConn
		nRead, addr, err :=
			readFromUDPConn(ob.conn, encryptedReq, ob.Config.ReplyTimeout)
		if err == errClosed {
			break
		}
		if err != nil {
			_ = ob.logError(0xEA288A, err)
			continue
		}
		recv, err := ob.Config.Cipher.Decrypt(encryptedReq[:nRead])
		if err != nil {
			_ = ob.logError(0xE9AF2F, err)
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
				_ = ob.logError(0xE69C60, err)
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
			_ = ob.logError(0xE06B58, err)
			continue
		}
		deadline := time.Now().Add(ob.Config.WriteTimeout)
		err = ob.conn.SetWriteDeadline(deadline)
		if err != nil {
			_ = ob.logError(0xE0AD06, err)
			continue
		}
		nWrit, err := ob.conn.WriteTo(encryptedReply, addr)
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

// Stop stops the Receiver from listening and
// receiving data by closing its connection.
func (ob *Receiver) Stop() {
	if ob.conn == nil {
		return
	}
	err := ob.conn.Close()
	if err != nil {
		_ = ob.logError(0xE9C2D1, err)
	}
	ob.conn = nil
} //                                                                        Stop

// -----------------------------------------------------------------------------
// # Packet Handlers

// receiveFragment handles a tagFragment packet sent by a Sender, and
// sends back a confirmation packet (tagConfirmation) to the Sender.
func (ob *Receiver) receiveFragment(recv []byte) ([]byte, error) {
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
		return nil, ob.logError(0xEF27F8, "bad 'sn':", sn)
	}
	index--
	packetCount, err := strconv.Atoi(count)
	if err != nil {
		return nil, ob.logError(0xE18A95, "bad 'count'")
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
		return nil, ob.logError(0xE1A99A, "unknown packet alteration")
	}
	if it.IsLoaded() {
		if ob.ReceiveData == nil {
			return nil, ob.logError(0xE49E2A, "nil Receiver.ReceiveData")
		}
		data, err := it.UnpackBytes(ob.Config.Compressor)
		if err != nil {
			return nil, ob.logError(0xE3DB1D, err)
		}
		err = ob.ReceiveData(it.Name, data)
		if err != nil {
			return nil, ob.logError(0xE77B4D, err)
		}
		ob.logInfo("received:", it.Name)
		if ob.Config.VerboseReceiver {
			it.LogStats("receiveFragment", ob.logInfo)
		}
		it.Reset()
	}
	confirmedHash := getHash(recv)
	reply := append([]byte(tagConfirmation), confirmedHash...)
	return reply, nil
} //                                                             receiveFragment

// sendDataItemHash handles a tagDataItemHash sent by a Sender.
func (ob *Receiver) sendDataItemHash(req []byte) ([]byte, error) {
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
	hash := getHash(data)
	reply := []byte(tagDataItemHash + fmt.Sprintf("%X", hash))
	return reply, nil
} //                                                            sendDataItemHash

// -----------------------------------------------------------------------------
// # Logging Methods

// logError returns a new error value generated by joining id and args
// and optionally calls Receiver.LogFunc (if not nil) to log the error.
func (ob *Receiver) logError(id uint32, args ...interface{}) error {
	ret := makeError(id, args...)
	if ob.Config != nil && ob.Config.LogFunc != nil {
		msg := ret.Error()
		ob.Config.LogFunc(msg)
	}
	return ret
} //                                                                    logError

// logInfo calls Receiver.LogFunc (if not nil) to log a message.
func (ob *Receiver) logInfo(args ...interface{}) {
	if ob.Config != nil && ob.Config.LogFunc != nil {
		ob.Config.LogFunc(args...)
	}
} //                                                                     logInfo

// end
