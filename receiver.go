// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[receiver.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// Receive(
//     port int,
//     cryptoKey []byte,
//     receive func(k string, v []byte) error,
// ) error
//
// type Receiver struct
//
// # Public Methods
//   ) Run() error
//   ) Stop()
//
// # Run() Internals
//   ) initRun() error
//   ) initRunDI(
//   ) buildReply(recv []byte) (reply []byte, err error)
//   ) sendReply(conn netUDPConn, addr net.Addr, reply []byte)
//
// # Packet Handlers
//   type fragmentHeader struct
//   ) readFragmentHeader(recv []byte) (*fragmentHeader, error)
//   ) receiveFragment(recv []byte) ([]byte, error)
//
// # Logging Methods
//   ) logError(id uint32, a ...interface{}) error
//   ) logInfo(a ...interface{})

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Receive sets up and runs a Receiver.
//
// Once it starts running, this function will only
// exit if the context is done or cancelled.
//
// It will only return an error if it fails to start
// because the port or cryptoKey is invalid.
//
func Receive(
	ctx context.Context,
	port int,
	cryptoKey []byte,
	receive func(k string, v []byte) error,
) error {
	ch := make(chan error, 1)
	var rc Receiver
	go func() {
		rc = Receiver{Port: port, CryptoKey: cryptoKey, Receive: receive}
		err := rc.Run()
		ch <- err
	}()
	defer rc.Stop()
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return nil
	}
} //                                                                     Receive

// -----------------------------------------------------------------------------

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

	// Receive is a callback function you must specify. This Receiver
	// will call it when a data item has been fully transferred.
	//
	// 'k' and 'v' will contain the key and value sent
	// by Sender.Send() or Sender.SendString(), etc.
	//
	// The reason there are two parameters is to separate metadata like
	// timestamps or filenames from the content of the transferred resource.
	//
	Receive func(k string, v []byte) error

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
// It calls Receive when a data transfer is complete, after the
// receiver has received, decrypted and re-assembled a data item.
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
// # Run() Internals

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
	if rc.Receive == nil {
		return rc.logError(0xE82C9E, "nil Receiver.Receive")
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
		rc.conn = nil // avoid non-nil interface with nil concrete value
		return rc.logError(0xEBF95F, err)
	}
	return nil
} //                                                                   initRunDI

// buildReply builds a reply to the received data. Presently, the only packet
// type received is a fragment (FRAG), replied with confirmation (CONF) packet.
func (rc *Receiver) buildReply(recv []byte) (reply []byte, err error) {
	switch {
	case len(recv) == 0:
		_ = rc.logError(0xE6B3BA, "received no data")
		//
	case bytes.HasPrefix(recv, []byte(tagFragment)):
		reply, err = rc.receiveFragment(recv)
		//
	default:
		reply = []byte("invalid_packet_header")
		err = rc.logError(0xE985CC, "invalid packet header")
	}
	return reply, err
} //                                                                  buildReply

// sendReply sends 'reply' to the specified connection
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
	h, err := rc.readFragmentHeader(recv)
	if err != nil {
		return nil, err
	}
	it := &rc.receivingDataItem
	it.Retain(h.key, h.hash, h.packetCount)
	compressedData := recv[h.dataOffset:]
	if len(compressedData) < 1 {
		return nil, rc.logError(0xE92B0F, "received no data")
	}
	// store the current piece
	if len(it.CompressedPieces[h.index]) == 0 {
		it.CompressedPieces[h.index] = compressedData
	} else if !bytes.Equal(compressedData, it.CompressedPieces[h.index]) {
		return nil, rc.logError(0xE1A99A, "unknown packet alteration")
	}
	if it.IsLoaded() {
		if rc.Receive == nil {
			return nil, rc.logError(0xE49E2A, "nil Receiver.Receive")
		}
		data, err := it.UnpackBytes(rc.Config.Compressor)
		if err != nil {
			return nil, rc.logError(0xE3DB1D, err)
		}
		err = rc.Receive(it.Key, data)
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

// -----------------------------------------------------------------------------
// # Logging Methods

// logError returns a new error generated by joining 'id' and 'a' and
// prints to Receiver.Config.LogWriter (if not nil) to log the error.
func (rc *Receiver) logError(id uint32, a ...interface{}) error {
	ret := makeError(id, a...)
	if rc.Config != nil && rc.Config.LogWriter != nil {
		s := ret.Error()
		rc.Config.LogWriter.Write([]byte(s))
	}
	return ret
} //                                                                    logError

// logInfo writes to Receiver.Config.LogWriter (if not nil) to log a message.
func (rc *Receiver) logInfo(a ...interface{}) {
	if rc.Config != nil && rc.Config.LogWriter != nil {
		fmt.Fprintln(rc.Config.LogWriter, a...)
	}
} //                                                                     logInfo

// end
