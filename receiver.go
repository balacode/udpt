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
	// must be shared by the Sender and the Receiver.
	//
	// The correct size of this key depends
	// on the implementation of SymmetricCipher.
	//
	CryptoKey []byte

	// Config contains UDP and other configuration settings.
	// These settings normally don't need to be changed.
	Config *Configuration

	// ReceiveData is a callback function you must specify. This Receiver
	// will call it when a data item has been fully transferred.
	//
	// 'k' and 'v' will contain the key and value sent
	// by Sender.Send() or Sender.SendString(), etc.
	//
	// The reason there are two parameters is to separate metadata like
	// timestamps or filenames from the content of the transferred resource.
	//
	ReceiveData func(k string, v []byte) error

	// ProvideData is a callback function you must specify. This
	// Receiver will call it to read back the named data item.
	//
	// This is needed to send back a confirmation hash to the Sender
	// to confirm the transfer. The Receiver carries out the hashing.
	//
	// If the resource specified by key 'k' is found,
	// it should return its bytes with a nil error.
	//
	// If the resource is not found, return (nil, nil).
	//
	// If there is some error, return (nil, <error>).
	//
	// NOTE: This callback may be removed because the hashing could be done
	// internally by the Receiver, so there'll be no need for a callback.
	//
	ProvideData func(k string) ([]byte, error)

	// -------------------------------------------------------------------------

	// conn is the UDP connection on which Receiver listens;
	// setting this to nil allows Run() to stop listening
	conn netUDPConn

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
func (rc *Receiver) Run() error {
	defer rc.Stop()
	if rc.Config == nil {
		rc.Config = NewDefaultConfig()
	}
	err := rc.initRun()
	if err != nil {
		return err
	}
	// receive transmissions
	encReq := make([]byte, rc.Config.PacketSizeLimit)
	for rc.conn != nil {
		// 'encReq' is overwritten after every readAndDecrypt
		recv, addr, err := readAndDecrypt(rc.conn, rc.Config.ReplyTimeout,
			rc.Config.Cipher, encReq)
		if err == errClosed {
			break
		}
		if err != nil {
			_ = rc.logError(0xEA288A, err)
			continue
		}
		if rc.Config.VerboseReceiver {
			rc.logInfo()
			rc.logInfo(strings.Repeat("-", 80))
			rc.logInfo("Receiver read", len(recv), "bytes from", addr)
		}
		reply, err := rc.buildReply(recv)
		if len(reply) == 0 || err != nil {
			continue
		}
		encReply, err := rc.Config.Cipher.Encrypt(reply)
		if err != nil {
			_ = rc.logError(0xE5C3E8, err)
			continue
		}
		rc.sendReply(rc.conn, addr, encReply)
	}
	return nil
} //                                                                         Run

// Stop stops the Receiver from listening and
// receiving data by closing its connection.
func (rc *Receiver) Stop() {
	if rc.conn == nil {
		return
	}
	err := rc.conn.Close()
	if err != nil {
		_ = rc.logError(0xE9C2D1, err)
	}
	rc.conn = nil
} //                                                                        Stop

// -----------------------------------------------------------------------------
// # Run() Helpers

// initRun checks if the receiver is properly configured
// and starts listening on the configured UDP address.
func (rc *Receiver) initRun() error {
	return rc.initRunDI(net.ResolveUDPAddr, net.ListenUDP)
} //                                                                     initRun

// initRunDI is only used by initRun() and provides parameters
// for dependency injection, to enable mocking during testing.
func (rc *Receiver) initRunDI(
	netResolveUDPAddr func(network string, addr string) (*net.UDPAddr, error),
	netListenUDP func(network string, laddr *net.UDPAddr) (*net.UDPConn, error),
) error {
	err := rc.Config.Validate()
	if err != nil {
		return rc.logError(0xE14BC8, err)
	}
	if rc.Port < 1 || rc.Port > 65535 {
		return rc.logError(0xE58B2F, "invalid Receiver.Port:", rc.Port)
	}
	err = rc.Config.Cipher.SetKey(rc.CryptoKey)
	if err != nil {
		return rc.logError(0xE8A5C6, "invalid Receiver.CryptoKey:", err)
	}
	if rc.ReceiveData == nil {
		return rc.logError(0xE82C9E, "nil Receiver.ReceiveData")
	}
	if rc.ProvideData == nil {
		return rc.logError(0xE48CC6, "nil Receiver.ProvideData")
	}
	udpAddr, err := netResolveUDPAddr("udp",
		fmt.Sprintf("0.0.0.0:%d", rc.Port))
	if err != nil {
		return rc.logError(0xE1D68C, err)
	}
	if rc.Config.VerboseReceiver {
		rc.logInfo(strings.Repeat("-", 80))
		rc.logInfo("Receiver listening...")
	}
	rc.conn, err = netListenUDP("udp", udpAddr)
	if err != nil {
		return rc.logError(0xEBF95F, err)
	}
	return nil
} //                                                                   initRunDI

func (rc *Receiver) buildReply(recv []byte) (reply []byte, err error) {
	switch {
	case len(recv) == 0:
		_ = rc.logError(0xE6B3BA, "received no data")
	case bytes.HasPrefix(recv, []byte(tagDataItemHash)):
		reply, err = rc.sendDataItemHash(recv)
	case bytes.HasPrefix(recv, []byte(tagFragment)):
		reply, err = rc.receiveFragment(recv)
	default:
		reply = []byte("invalid_packet_header")
		err = rc.logError(0xE985CC, "invalid packet header")
	}
	return reply, err
} //                                                                  buildReply

func (rc *Receiver) sendReply(conn netUDPConn, addr net.Addr, reply []byte) {
	deadline := time.Now().Add(rc.Config.WriteTimeout)
	err := conn.SetWriteDeadline(deadline)
	if err != nil {
		_ = rc.logError(0xE0AD06, err)
		return
	}
	nWrit, err := conn.WriteTo(reply, addr)
	if err != nil {
		_ = rc.logError(0xEA63C4, err)
		return
	}
	if rc.Config.VerboseReceiver {
		rc.logInfo("Receiver wrote", nWrit, "bytes to", addr)
	}
} //                                                                   sendReply

// -----------------------------------------------------------------------------
// # Packet Handlers

// fragmentHeader contains details read from a received fragment
type fragmentHeader struct {
	dataOffset  int    // position of compressed data (part of the value)
	key         string // key 'k' of the key-value message
	hash        []byte // hash of entire key-value message
	index       int    // 0-based index of this fragment
	packetCount int    // total number of fragments (i.e. packets) in message
}

// readFragmentHeader reads the header from a received fragment packet
func (rc *Receiver) readFragmentHeader(recv []byte) (*fragmentHeader, error) {
	if !bytes.HasPrefix(recv, []byte(tagFragment)) {
		return nil, rc.logError(0xE4F3C5, "missing header")
	}
	var h fragmentHeader
	h.dataOffset = bytes.Index(recv, []byte("\n"))
	if h.dataOffset == -1 {
		return nil, rc.logError(0xE6CF52, "newline not found")
	}
	h.dataOffset++ // skip newline
	//
	s := string(recv[len(tagFragment):h.dataOffset])
	h.key = getPart(s, "key:", " ")
	//
	var err error
	h.hash, err = hex.DecodeString(getPart(s, "hash:", " "))
	if err != nil || len(h.hash) != 32 {
		return nil, rc.logError(0xEB6CB7, "bad hash")
	}
	h.packetCount, _ = strconv.Atoi(getPart(s, "count:", "\n"))
	if h.packetCount < 1 {
		return nil, rc.logError(0xE18A95, "bad 'count'")
	}
	h.index, _ = strconv.Atoi(getPart(s, "sn:", " "))
	if h.index < 1 || h.index > h.packetCount {
		return nil, rc.logError(0xEF27F8, "bad 'sn'")
	}
	h.index--
	return &h, nil
} //                                                          readFragmentHeader

// receiveFragment handles a tagFragment packet sent by a Sender, and
// sends back a confirmation packet (tagConfirmation) to the Sender.
func (rc *Receiver) receiveFragment(recv []byte) ([]byte, error) {
	if !bytes.HasPrefix(recv, []byte(tagFragment)) {
		return nil, rc.logError(0xE4F3C5, "missing header")
	}
	err := rc.Config.Validate()
	if err != nil {
		return nil, rc.logError(0xE9B5C7, err)
	}
	dataOffset := bytes.Index(recv, []byte("\n"))
	if dataOffset == -1 {
		return nil, rc.logError(0xE6CF52, "newline not found")
	}
	dataOffset++ // skip newline
	var (
		header  = string(recv[len(tagFragment):dataOffset])
		key     = getPart(header, "key:", " ")
		hexHash = getPart(header, "hash:", " ")
		sn      = getPart(header, "sn:", " ")
		count   = getPart(header, "count:", "\n")
	)
	packetCount, err := strconv.Atoi(count)
	if err != nil || packetCount < 1 {
		return nil, rc.logError(0xE18A95, "bad 'count'")
	}
	index, err := strconv.Atoi(sn)
	if err != nil {
		return nil, rc.logError(0xEF27F8, "bad 'sn':", sn)
	}
	if index < 1 || index > packetCount {
		return nil, rc.logError(0xE8CF4D,
			"'sn'", index, "out of range 1 -", packetCount)
	}
	index--
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return nil, rc.logError(0xE5CA62, err)
	}
	if len(hash) != 32 {
		return nil, rc.logError(0xEB6CB7, "bad hash size")
	}
	it := &rc.receivingDataItem
	it.Retain(key, hash, packetCount)
	compressedData := recv[dataOffset:]
	if len(compressedData) < 1 {
		return nil, rc.logError(0xE92B0F, "received no data")
	}
	// store the current piece
	if len(it.CompressedPieces[index]) == 0 {
		it.CompressedPieces[index] = compressedData
	} else if !bytes.Equal(compressedData, it.CompressedPieces[index]) {
		return nil, rc.logError(0xE1A99A, "unknown packet alteration")
	}
	if it.IsLoaded() {
		if rc.ReceiveData == nil {
			return nil, rc.logError(0xE49E2A, "nil Receiver.ReceiveData")
		}
		data, err := it.UnpackBytes(rc.Config.Compressor)
		if err != nil {
			return nil, rc.logError(0xE3DB1D, err)
		}
		err = rc.ReceiveData(it.Key, data)
		if err != nil {
			return nil, rc.logError(0xE77B4D, err)
		}
		rc.logInfo("received:", it.Key)
		if rc.Config.VerboseReceiver {
			var sb strings.Builder
			it.LogStats("receiveFragment", &sb)
			rc.logInfo(sb.String())
		}
		it.Reset()
	}
	confirmedHash := getHash(recv)
	reply := append([]byte(tagConfirmation), confirmedHash...)
	return reply, nil
} //                                                             receiveFragment

// sendDataItemHash handles a tagDataItemHash sent by a Sender.
func (rc *Receiver) sendDataItemHash(req []byte) ([]byte, error) {
	if !bytes.HasPrefix(req, []byte(tagDataItemHash)) {
		return nil, rc.logError(0xE7B653, "missing header")
	}
	if rc.ProvideData == nil {
		return nil, rc.logError(0xE73A1C, "nil ProvideData")
	}
	k := string(req[len(tagDataItemHash):])
	v, err := rc.ProvideData(k)
	if err != nil {
		return nil, rc.logError(0xE7F7C9, err)
	}
	hash := getHash(v)
	reply := []byte(tagDataItemHash + fmt.Sprintf("%X", hash))
	return reply, nil
} //                                                            sendDataItemHash

// -----------------------------------------------------------------------------
// # Logging Methods

// logError returns a new error generated by joining 'id' and 'a' and
// optionally calls Receiver.Config.LogFunc (if not nil) to log the error.
func (rc *Receiver) logError(id uint32, a ...interface{}) error {
	ret := makeError(id, a...)
	if rc.Config != nil && rc.Config.LogFunc != nil {
		msg := ret.Error()
		rc.Config.LogFunc(msg)
	}
	return ret
} //                                                                    logError

// logInfo calls Receiver.Config.LogFunc (if not nil) to log a message.
func (rc *Receiver) logInfo(a ...interface{}) {
	if rc.Config != nil && rc.Config.LogFunc != nil {
		rc.Config.LogFunc(a...)
	}
} //                                                                     logInfo

// end
