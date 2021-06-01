// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                         /[sender.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// # Functions
//
//   Send(addr string, k string, v []byte, cryptoKey []byte,
//       config ...*Configuration,
//   ) error
//
//   SendString(addr string, k, v string, cryptoKey []byte,
//       config ...*Configuration,
//   ) error
//
// # Sender Type
//   Sender struct
//
// # Main Methods (sd *Sender)
//   ) Send(k string, v []byte) error
//   ) SendString(k, v string) error
//
// # Informatory Properties (sd *Sender)
//   ) AverageResponseMs() float64
//   ) DeliveredAllParts() bool
//   ) TransferSpeedKBpS() float64
//
// # Informatory Methods (sd *Sender)
//   ) LogStats(w ...io.Writer)
//
// # Internal Lifecycle Methods (sd *Sender)
//   ) beginSend(k string, v []byte) (hash []byte, err error)
//   ) makePackets(k string, comp []byte) error
//   ) connect() (netUDPConn, error)
//   ) connectDI( . . .
//   ) sendUndeliveredPackets() error
//   ) collectConfirmations()
//   ) waitForAllConfirmations()
//   ) close()
//   ) endSend(k string, hash []byte) error
//
// # Internal Helper Methods (sd *Sender)
//   ) logError(id uint32, a ...interface{}) error
//   ) logInfo(a ...interface{})
//   ) makePacket(data []byte) (*senderPacket, error)
//   ) validateAddress() error

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------
// # Functions

// Send creates a Sender and uses it to transfer a key-value
// pair to the Receiver specified by address 'addr'.
//
// addr specifies the host and port number of the Receiver,
// for example "website.com:9876" or "127.0.0.1:9876"
//
// k is any string you want to use as the key. It can be blank if not needed.
// It could be a filename, timestamp, UUID, or some other metadata that
// gives context to the value being sent.
//
// v is the value being sent as a sequence of bytes. It can be as large
// as the free memory available on the Sender's and Receiver's machine.
//
// cryptoKey is the symmetric encryption key shared by the
// Sender and Receiver and used to encrypt the sent message.
//
// config is an optional Configuration you can customize. If you leave it out,
// Send() will use the configuration returned by NewDefaultConfig().
//
func Send(addr string, k string, v []byte, cryptoKey []byte,
	config ...*Configuration,
) error {
	if len(config) > 1 {
		return makeError(0xE8C0D4, "too many 'config' arguments")
	}
	var cf *Configuration
	if len(config) == 1 {
		cf = config[0]
	}
	if cf == nil {
		cf = NewDefaultConfig()
	}
	sender := Sender{Address: addr, CryptoKey: cryptoKey, Config: cf}
	err := sender.Send(k, v)
	return err
} //                                                                        Send

// SendString creates a Sender and uses it to transfer a key-value
// pair of strings to the Receiver specified by address 'addr'.
//
// addr specifies the host and port number of the Receiver,
// for example "website.com:9876" or "127.0.0.1:9876"
//
// k is any string you want to use as the key. It can be blank if not needed.
// It could be a filename, timestamp, UUID, or some other metadata that
// gives context to the value being sent.
//
// v is the value being sent as a string. It can be as large as the
// free memory available on the Sender's and Receiver's machine.
//
// cryptoKey is the symmetric encryption key shared by the Sender
// and Receiver and used to encrypt the sent message.
//
// config is an optional Configuration you can customize. If you leave it out,
// SendString() will use the configuration returned by NewDefaultConfig().
//
func SendString(addr string, k, v string, cryptoKey []byte,
	config ...*Configuration,
) error {
	return Send(addr, k, []byte(v), cryptoKey, config...)
} //                                                                  SendString

// -----------------------------------------------------------------------------
// # Sender Type

// udpStats contains UDP transfer statistics, such as the transfer
// speed and the number of packets delivered and lost.
type udpStats struct {
	bytesDelivered   int64
	bytesLost        int64
	packetsDelivered int64
	packetsLost      int64
	transferTime     time.Duration
} //                                                                    udpStats

// Sender coordinates sending key-value messages to a listening Receiver.
//
// You can use standalone Send() and SendString() functions
// to create a single-use Sender to send a message, but it's
// more efficient to construct a reusable Sender.
//
type Sender struct {

	// Address is the domain name or IP address of the listening
	// receiver with the port number. For example: "127.0.0.1:9876"
	//
	// The port number must be between 1 and 65535.
	//
	Address string

	// CryptoKey is the secret symmetric encryption key that
	// must be shared by the Sender and the Receiver.
	//
	// The correct size of this key depends on
	// the implementation of SymmetricCipher.
	//
	CryptoKey []byte

	// Config contains UDP and other configuration settings.
	// These settings normally don't need to be changed.
	Config *Configuration

	// -------------------------------------------------------------------------

	// conn holds the UDP connection to a Receiver
	conn netUDPConn

	// dataHash contains the hash of all bytes of the data item being sent
	dataHash []byte

	// stats contains UDP transfer statistics, such as the transfer
	// speed and the number of packets delivered and lost
	stats udpStats

	// packets contains all the packets of the currently transferred data item;
	// some of them may have been delivered, while others may need (re)sending
	packets []senderPacket

	// startTime is the time the first packet was sent, after
	// the bytes of the data item have been compressed
	startTime time.Time

	// wg is used by waitForAllConfirmations()
	// to wait for sendUndeliveredPackets()
	wg sync.WaitGroup
} //                                                                      Sender

// -----------------------------------------------------------------------------
// # Main Methods (sd *Sender)

// Send transfers a key-value to the Receiver specified by Sender.Address.
//
// 'k' is any string you want to use as the key. It can be blank if not needed.
// It could be a filename, timestamp, UUID, or some other metadata that
// gives context to the value being sent.
//
// 'v' is the value being sent as a sequence of bytes. It can be as large
// as the free memory available on the Sender's and Receiver's machine.
//
func (sd *Sender) Send(k string, v []byte) error {
	return sd.sendDI(k, v, sd.connect, sd.sendUndeliveredPackets)
} //                                                                        Send

// sendDI is only used by Send() and provides parameters for
// dependency injection, to enable mocking during testing.
func (sd *Sender) sendDI(k string, v []byte,
	connect func() (netUDPConn, error),
	sendUndeliveredPackets func() error,
) error {
	if sd.Config == nil {
		sd.Config = NewDefaultConfig()
	}
	hash, err := sd.beginSend(k, v)
	if hash == nil {
		return err
	}
	newConn, err := connect()
	if err != nil {
		return sd.logError(0xE8B8D0, err)
	}
	sd.conn = newConn
	go sd.collectConfirmations() // exits when conn becomes nil
	for retries := 0; retries < sd.Config.SendRetries; retries++ {
		err = sendUndeliveredPackets()
		if err != nil {
			defer func() { sd.close() }()
			return sd.logError(0xE23CE0, err)
		}
		sd.waitForAllConfirmations()
		if sd.DeliveredAllParts() {
			break
		}
		time.Sleep(sd.Config.SendRetryInterval)
	}
	sd.close()
	return sd.endSend(k, hash)
} //                                                                      sendDI

// SendString transfers a key and value string
// to the Receiver specified by Sender.Address.
//
// 'k' is any string you want to use as the key. It can be blank if not needed.
// It could be a filename, timestamp, UUID, or some other metadata that
// gives context to the value being sent.
//
// 'v' is the value being sent as a string. It can be as large as the
// free memory available on the Sender's and Receiver's machine.
//
func (sd *Sender) SendString(k string, v string) error {
	return sd.Send(k, []byte(v))
} //                                                                  SendString

// -----------------------------------------------------------------------------
// # Informatory Properties (sd *Sender)

// AverageResponseMs is the average response time, in milliseconds, between
// a packet being sent and its delivery confirmation being received.
func (sd *Sender) AverageResponseMs() float64 {
	if sd.stats.packetsDelivered == 0 {
		return 0.0
	}
	// instead of using transferTime.Milliseconds(),
	// cast to float64 to get sub-millisecond timing
	ret := float64(sd.stats.transferTime) /
		float64(time.Millisecond) /
		float64(sd.stats.packetsDelivered)
	return ret
} //                                                           AverageResponseMs

// DeliveredAllParts returns true if all parts of the
// sent data item have been delivered. I.e. all packets
// have been sent, resent if needed, and confirmed.
func (sd *Sender) DeliveredAllParts() bool {
	defer func() {
		if r := recover(); r != nil {
			_ = sd.logError(0xEC7A22, "Sender.DeliveredAllParts panic:", r)
		}
	}()
	ret := len(sd.packets) > 0
	for _, pk := range sd.packets {
		if !bytes.Equal(pk.sentHash, pk.confirmedHash) {
			ret = false
			break
		}
	}
	return ret
} //                                                           DeliveredAllParts

// TransferSpeedKBpS returns the transfer speed of the current Send
// operation, in Kilobytes (more accurately, Kibibytes) per second.
func (sd *Sender) TransferSpeedKBpS() float64 {
	if sd.stats.transferTime < 1 {
		return 0.0
	}
	sec := float64(sd.stats.transferTime) / float64(time.Second)
	ret := float64(sd.stats.bytesDelivered/1024) / sec
	return ret
} //                                                           TransferSpeedKBpS

// -----------------------------------------------------------------------------
// # Informatory Methods (sd *Sender)

// LogStats prints UDP transfer statistics to the specified writer 'wr',
// or to Sender.LogWriter. If both are not specified, does nothing.
func (sd *Sender) LogStats(w ...io.Writer) {
	//
	log := sd.logInfo
	if len(w) > 0 {
		log = func(a ...interface{}) { fmt.Fprintln(w[0], a...) }
	}
	tItem := time.Duration(0)
	for i, pk := range sd.packets {
		tPacket, status := time.Duration(0), "✔"
		if pk.IsDelivered() {
			if !pk.confirmedTime.IsZero() {
				tPacket = pk.confirmedTime.Sub(pk.sentTime)
			}
		} else {
			status = "LOST"
		}
		var (
			sn = padf(4, "%d", i)
			t0 = pk.sentTime.String()[:24]
			t1 = pk.confirmedTime.String()[:24]
			ms = fmt.Sprintf("%0.1f ms",
				float64(tPacket)/float64(time.Millisecond))
		)
		if pk.confirmedTime.IsZero() {
			t1 = "NONE"
		}
		log("SN:", sn, "T0:", t0, "T1:", t1, status, ms)
		tItem += tPacket
	}
	var (
		sec   = sd.stats.transferTime.Seconds()
		avg   = sd.AverageResponseMs()
		speed = sd.TransferSpeedKBpS()
		prt   = func(tag, format string, v interface{}) {
			log(tag, fmt.Sprintf(format, v))
		}
	)
	prt("B. delivered:", "%d", sd.stats.bytesDelivered)
	prt("Bytes lost  :", "%d", sd.stats.bytesLost)
	prt("P. delivered:", "%d", sd.stats.packetsDelivered)
	prt("Packets lost:", "%d", sd.stats.packetsLost)
	prt("Time in item:", "%0.1f s", sec)
	prt("Avg./ Packet:", "%0.1f ms", avg)
	prt("Trans. speed:", "%0.1f KiB/s", speed)
} //                                                                    LogStats

// -----------------------------------------------------------------------------
// # Internal Lifecycle Methods (sd *Sender)

// beginSend checks if the sender is properly configured before sending
func (sd *Sender) beginSend(k string, v []byte) (hash []byte, err error) {
	//
	// setup cipher
	if sd.Config.Cipher == nil {
		return nil, sd.logError(0xE83D07, "nil Sender.Config.Cipher")
	}
	err = sd.Config.Cipher.SetKey(sd.CryptoKey)
	if err != nil {
		return nil, sd.logError(0xE02D7B, "invalid Sender.CryptoKey:", err)
	}
	// check settings
	err = sd.Config.Validate()
	if err != nil {
		return nil, sd.logError(0xE5D92D, "invalid Sender.Config:", err)
	}
	err = sd.validateAddress()
	if err != nil {
		return nil, sd.logError(0xE5A04A, err)
	}
	hash = getHash(v)
	if sd.Config.VerboseSender {
		sd.logInfo("\n" + strings.Repeat("-", 80) + "\n" +
			fmt.Sprintf("Send key: %s size: %d hash: %X",
				k, len(v), hash))
	}
	comp, err := sd.Config.Compressor.Compress(v)
	if err != nil {
		return nil, sd.logError(0xE2EB59, err)
	}
	sd.dataHash = hash
	sd.startTime = time.Now()
	err = sd.makePackets(k, comp)
	if err != nil {
		return nil, err
	}
	return hash, nil
} //                                                                   beginSend

// makePackets creates the packets for sending over UDP,
// by partitioning compressed message 'comp'
func (sd *Sender) makePackets(k string, comp []byte) error {
	length := len(comp)
	if length == 0 {
		sd.packets = nil
		return nil
	}
	max := sd.Config.PacketPayloadSize
	n := length / max
	if (n * max) < length {
		n++
	}
	packets := make([]senderPacket, n)
	for i := range packets {
		a := i * max
		b := a + max
		if b > len(comp) {
			b = len(comp)
		}
		header := tagFragment + fmt.Sprintf(
			"key:%s hash:%X sn:%d count:%d\n",
			k, sd.dataHash, i+1, n,
		)
		pk, err := sd.makePacket(append([]byte(header), comp[a:b]...))
		if err != nil {
			return sd.logError(0xE567A4, err)
		}
		packets[i] = *pk
	}
	sd.packets = packets
	return nil
} //                                                                 makePackets

// connect connects to the Receiver at Sender.Address and
// returns a new UDP connection or nil and an error instance.
//
// Note that it doesn't change the value of Sender.conn
//
func (sd *Sender) connect() (netUDPConn, error) {
	fn := func(network string, laddr *net.UDPAddr, raddr *net.UDPAddr,
	) (netUDPConn, error) {
		return net.DialUDP(network, laddr, raddr)
	}
	return sd.connectDI(fn)
} //                                                                     connect

// connectDI is only used by connect() and provides a parameter
// for dependency injection, to enable mocking during testing.
func (sd *Sender) connectDI(
	netDialUDP func(string, *net.UDPAddr, *net.UDPAddr) (netUDPConn, error),
) (netUDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", sd.Address)
	if err != nil {
		return nil, sd.logError(0xEC7C6B, "ResolveUDPAddr:", err)
	}
	var conn netUDPConn
	conn, err = netDialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, sd.logError(0xE15CE1, err)
	}
	err = conn.SetWriteBuffer(sd.Config.SendBufferSize)
	if err != nil {
		return nil, sd.logError(0xE5F9C7, err)
	}
	return conn, nil
} //                                                                   connectDI

// sendUndeliveredPackets sends all undelivered
// packets to the destination Receiver.
func (sd *Sender) sendUndeliveredPackets() error {
	n := len(sd.packets)
	for i := 0; i < n; i++ {
		pk := &sd.packets[i]
		if pk.IsDelivered() {
			continue
		}
		time.Sleep(sd.Config.SendPacketInterval)
		sd.wg.Add(1)
		go func() {
			err := pk.Send(sd.conn, sd.Config.Cipher)
			if err != nil {
				_ = sd.logError(0xE67BA4, err)
			}
			sd.wg.Done()
		}()
	}
	return nil
} //                                                      sendUndeliveredPackets

// collectConfirmations enters a loop that receives confirmation packets
// from the sender, and marks all confirmed packets as delivered.
func (sd *Sender) collectConfirmations() {
	encReply := make([]byte, sd.Config.PacketSizeLimit)
	for sd.conn != nil {
		// 'encReply' is overwritten after every readAndDecrypt
		recv, addr, err := readAndDecrypt(sd.conn, sd.Config.ReplyTimeout,
			sd.Config.Cipher, encReply)
		if err == errClosed {
			break
		}
		if err != nil {
			_ = sd.logError(0xE9D1CC, err)
			continue
		}
		if !bytes.HasPrefix(recv, []byte(tagConfirmation)) {
			_ = sd.logError(0xE96D3B, "bad reply header")
			if sd.Config.VerboseSender {
				sd.logInfo("ERROR received:", len(recv), "bytes")
			}
			continue
		}
		confirmedHash := recv[len(tagConfirmation):]
		if sd.Config.VerboseSender {
			sd.logInfo("Sender received", len(recv), "bytes from", addr)
		}
		go func(confirmedHash []byte) {
			for i, pk := range sd.packets {
				if bytes.Equal(pk.sentHash, confirmedHash) {
					sd.packets[i].confirmedTime = time.Now()
					sd.packets[i].confirmedHash = confirmedHash
					break
				}
			}
		}(confirmedHash)
	}
} //                                                        collectConfirmations

// waitForAllConfirmations waits for all confirmation packets to
// be received from the receiver. Since UDP packet delivery is not
// guaranteed, some confirmations may not be received. This method
// will only wait for the duration specified in Config.ReplyTimeout.
func (sd *Sender) waitForAllConfirmations() {
	if sd.Config.VerboseSender {
		sd.logInfo("Waiting . . .")
	}
	t0 := time.Now()
	sd.wg.Wait()
	for {
		time.Sleep(sd.Config.SendWaitInterval)
		if sd.DeliveredAllParts() {
			if sd.Config.VerboseSender {
				sd.logInfo("Delivered all packets")
			}
			break
		}
		since := time.Since(t0)
		if since >= sd.Config.ReplyTimeout {
			sd.logInfo("Config.ReplyTimeout exceeded",
				fmt.Sprintf("%0.1f", since.Seconds()))
			break
		}
	}
	for _, pk := range sd.packets {
		if pk.IsDelivered() {
			sd.stats.bytesDelivered += int64(len(pk.data))
			sd.stats.packetsDelivered++
		} else {
			sd.stats.bytesLost += int64(len(pk.data))
			sd.stats.packetsLost++
		}
	}
	if sd.Config.VerboseSender {
		sd.logInfo("Waited:", time.Since(t0))
	}
} //                                                     waitForAllConfirmations

// close closes the UDP connection.
func (sd *Sender) close() {
	if sd.conn == nil {
		return
	}
	err := sd.conn.Close()
	sd.conn = nil
	if err != nil {
		_ = sd.logError(0xEA7D7E, err)
	}
} //                                                                       close

// endSend finializes the Send() by checking if the message was successfully
// delivered by requesting a confirmation from the Receiver.
func (sd *Sender) endSend(k string, hash []byte) error {
	if !sd.DeliveredAllParts() {
		return sd.logError(0xE1C3A7, "undelivered packets")
	}
	if sd.Config.VerboseSender {
		sd.LogStats()
	}
	return nil
} //                                                                     endSend

// -----------------------------------------------------------------------------
// # Internal Helper Methods (sd *Sender)

// logError returns a new error generated by joining 'id' and 'a' and
// prints to Sender.Config.LogWriter (if not nil) to log the error.
func (sd *Sender) logError(id uint32, a ...interface{}) error {
	ret := makeError(id, a...)
	if sd.Config != nil && sd.Config.LogWriter != nil {
		s := ret.Error()
		fmt.Fprintln(sd.Config.LogWriter, s)
	}
	return ret
} //                                                                    logError

// logInfo writes to Sender.Config.LogWriter (if not nil) to log a message.
func (sd *Sender) logInfo(a ...interface{}) {
	if sd.Config != nil && sd.Config.LogWriter != nil {
		fmt.Fprintln(sd.Config.LogWriter, a...)
	}
} //                                                                     logInfo

// makePacket prepares a packet for immediate sending: it stores,
// hashes data and sets the packet's sentTime to current time.
//
// The size of the packet must not exceed Config.PacketSizeLimit
//
func (sd *Sender) makePacket(data []byte) (*senderPacket, error) {
	if len(data) > sd.Config.PacketSizeLimit {
		return nil, sd.logError(0xE71F9B, "len(data) > Config.PacketSizeLimit")
	}
	sentHash := getHash(data)
	pk := senderPacket{
		data:     data,
		sentHash: sentHash,
		sentTime: time.Now(),
		// confirmedHash, confirmedTime: zero value
	}
	return &pk, nil
} //                                                                  makePacket

// validateAddress returns nil if Address is valid, or an error otherwise.
// Presently it only checks if the address contains a valid port number.
func (sd *Sender) validateAddress() error {
	ad := sd.Address
	if strings.TrimSpace(ad) == "" {
		return errors.New("missing Sender.Address")
	}
	var port int
	if i := strings.Index(ad, ":"); i != -1 {
		port, _ = strconv.Atoi(ad[i+1:])
	}
	if port < 1 || port > 65535 {
		return errors.New("invalid port in Sender.Address")
	}
	return nil
} //                                                             validateAddress

// end
