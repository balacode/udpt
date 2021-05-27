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
//   ) LogStats()
//
// # Internal Lifecycle Methods (sd *Sender)
//   ) requestDataItemHash(k string) []byte
//   ) connect() (netUDPConn, error)
//   ) connectDI( . . .
//   ) sendUndeliveredPackets() error
//   ) collectConfirmations()
//   ) waitForAllConfirmations()
//   ) close()
//
// # Internal Helper Methods (sd *Sender)
//   ) getPacketCount(length int) int
//   ) logError(id uint32, a ...interface{}) error
//   ) logInfo(a ...interface{})
//   ) makePacket(data []byte) (*senderPacket, error)

import (
	"bytes"
	"encoding/hex"
	"fmt"
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
// cryptoKey is the symmetric encryption key shared by the Sender
// and Receiver and used to encrypt the sent message.
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
	packsLost        int64
	transferTime     time.Duration
} //                                                                    udpStats

// Sender coordinates sending key-value messages to a listening Receiver.
//
// You can use standalone Send() and SendString() functions to create a
// single-use Sender to send a message, but it's more efficient to
// construct a reusable Sender.
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
	if sd.Config == nil {
		sd.Config = NewDefaultConfig()
	}
	// setup cipher
	if sd.Config.Cipher == nil {
		return sd.logError(0xE83D07, "nil Sender.Config.Cipher")
	}
	err := sd.Config.Cipher.SetKey(sd.CryptoKey)
	if err != nil {
		return sd.logError(0xE02D7B, "invalid Sender.CryptoKey:", err)
	}
	// check settings
	err = sd.Config.Validate()
	if err != nil {
		return sd.logError(0xE5D92D, "Invalid Sender.Config:", err)
	}
	if strings.TrimSpace(sd.Address) == "" {
		return sd.logError(0xE5A04A, "missing Sender.Address")
	}
	var port int
	if i := strings.Index(sd.Address, ":"); i != -1 {
		port, _ = strconv.Atoi(sd.Address[i+1:])
	}
	if port < 1 || port > 65535 {
		return sd.logError(0xE20BB9, "invalid port in Sender.Address")
	}
	// prepare for transfer
	hash := getHash(v)
	if sd.Config.VerboseSender {
		sd.logInfo("\n" + strings.Repeat("-", 80) + "\n" +
			fmt.Sprintf("Send key: %s size: %d hash: %X",
				k, len(v), hash))
	}
	remoteHash := sd.requestDataItemHash(k)
	if bytes.Equal(hash, remoteHash) {
		return nil
	}
	comp, err := sd.Config.Compressor.Compress(v)
	if err != nil {
		return sd.logError(0xE2EB59, err)
	}
	packetCount := sd.getPacketCount(len(comp))
	sd.dataHash = hash
	sd.startTime = time.Now()
	sd.packets = make([]senderPacket, packetCount)
	//
	// begin transfer
	for i := range sd.packets {
		a := i * sd.Config.PacketPayloadSize
		b := a + sd.Config.PacketPayloadSize
		if b > len(comp) {
			b = len(comp)
		}
		header := tagFragment + fmt.Sprintf(
			"key:%s hash:%X sn:%d count:%d\n",
			k, sd.dataHash, i+1, packetCount,
		)
		packet, err2 := sd.makePacket(
			append([]byte(header), comp[a:b]...),
		)
		if err2 != nil {
			return sd.logError(0xE567A4, err2)
		}
		sd.packets[i] = *packet
	}
	newConn, err := sd.connect()
	if err != nil {
		return sd.logError(0xE8B8D0, err)
	}
	sd.conn = newConn
	go sd.collectConfirmations() // exits when conn becomes nil
	for retries := 0; retries < sd.Config.SendRetries; retries++ {
		err = sd.sendUndeliveredPackets()
		if err != nil {
			defer func() {
				sd.close()
			}()
			return sd.logError(0xE23CE0, err)
		}
		sd.waitForAllConfirmations()
		if sd.DeliveredAllParts() {
			break
		}
		time.Sleep(sd.Config.SendRetryInterval)
	}
	sd.close()
	if !sd.DeliveredAllParts() {
		return sd.logError(0xE1C3A7, "undelivered packets")
	}
	remoteHash = sd.requestDataItemHash(k)
	if !bytes.Equal(hash, remoteHash) {
		return sd.logError(0xE1F101, "hash mismatch")
	}
	if sd.Config.VerboseSender {
		sd.LogStats()
	}
	return nil
} //                                                                        Send

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
	ret := true
	for _, packet := range sd.packets {
		if !bytes.Equal(packet.sentHash, packet.confirmedHash) {
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

// LogStats prints UDP transfer statistics using the passed logFunc function.
//
// logFunc should have a signature matching log.Println or fmt.Println.
// It is optional. If you omit it, uses Sender.Config.LogFunc for output.
//
// like log.Println: func(...interface{})
//
// like fmt.Println: func(...interface{}) (int, error)
//
func (sd *Sender) LogStats(logFunc ...interface{}) {
	//
	log := sd.logInfo // func(v ...interface{})
	if len(logFunc) > 0 {
		switch fn := logFunc[0].(type) {
		case func(...interface{}): // like log.Println
			log = fn
		case func(...interface{}) (int, error): // like fmt.Println
			log = func(v ...interface{}) { _, _ = fn(v...) }
		}
	}
	tItem := time.Duration(0)
	for i, pack := range sd.packets {
		tPack, status := time.Duration(0), "✔"
		if pack.IsDelivered() {
			if !pack.confirmedTime.IsZero() {
				tPack = pack.confirmedTime.Sub(pack.sentTime)
			}
		} else {
			status = "LOST"
		}
		var (
			sn = padf(4, "%d", i)
			t0 = pack.sentTime.String()[:24]
			t1 = pack.confirmedTime.String()[:24]
			ms = fmt.Sprintf("%0.1f ms",
				float64(tPack)/float64(time.Millisecond))
		)
		if pack.confirmedTime.IsZero() {
			t1 = "NONE"
		}
		log("SN:", sn, "T0:", t0, "T1:", t1, status, ms)
		tItem += tPack
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
	prt("Packets lost:", "%d", sd.stats.packsLost)
	prt("Time in item:", "%0.1f s", sec)
	prt("Avg./ Packet:", "%0.1f ms", avg)
	prt("Trans. speed:", "%0.1f KiB/s", speed)
} //                                                                    LogStats

// -----------------------------------------------------------------------------
// # Internal Lifecycle Methods (sd *Sender)

// requestDataItemHash requests and waits for the listening receiver to
// return the hash of the data item identified by key 'k'. If the receiver
// can locate the data item, it returns its hash, otherwise it returns nil.
func (sd *Sender) requestDataItemHash(k string) []byte {
	tempConn, err := sd.connect()
	if err != nil {
		_ = sd.logError(0xE7DF8B, err)
		return nil
	}
	packet, err := sd.makePacket([]byte(tagDataItemHash + k))
	if err != nil {
		_ = sd.logError(0xE34A8E, err)
		return nil
	}
	err = packet.Send(tempConn, sd.Config.Cipher)
	if err != nil {
		_ = sd.logError(0xE89B11, err)
		return nil
	}
	encReply := make([]byte, sd.Config.PacketSizeLimit)
	reply, _, err := readAndDecrypt(tempConn, sd.Config.ReplyTimeout,
		sd.Config.Cipher, encReply)
	if err != nil {
		_ = sd.logError(0xE97FC3, err)
		return nil
	}
	var hash []byte
	if len(reply) > 0 {
		if !bytes.HasPrefix(reply, []byte(tagDataItemHash)) {
			_ = sd.logError(0xE08AD4, "invalid reply:", reply)
			return nil
		}
		hexHash := string(reply[len(tagDataItemHash):])
		if hexHash == "not_found" {
			return nil
		}
		hash, err = hex.DecodeString(hexHash)
		if err != nil {
			_ = sd.logError(0xE6E7A9, err)
			return nil
		}
	}
	return hash
} //                                                         requestDataItemHash

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
		packet := &sd.packets[i]
		if packet.IsDelivered() {
			continue
		}
		time.Sleep(sd.Config.SendPacketInterval)
		sd.wg.Add(1)
		go func() {
			err := packet.Send(sd.conn, sd.Config.Cipher)
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
			for i, packet := range sd.packets {
				if bytes.Equal(packet.sentHash, confirmedHash) {
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
	sd.logInfo("Waiting . . .")
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
	for _, packet := range sd.packets {
		if packet.IsDelivered() {
			sd.stats.bytesDelivered += int64(len(packet.data))
			sd.stats.packetsDelivered++
		} else {
			sd.stats.bytesLost += int64(len(packet.data))
			sd.stats.packsLost++
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

// -----------------------------------------------------------------------------
// # Internal Helper Methods (sd *Sender)

// getPacketCount calculates the number of packets needed to send 'length'
// bytes. This depends on the setting of Config.PacketPayloadSize.
func (sd *Sender) getPacketCount(length int) int {
	if length < 1 {
		return 0
	}
	count := length / sd.Config.PacketPayloadSize
	if (count * sd.Config.PacketPayloadSize) < length {
		count++
	}
	return count
} //                                                              getPacketCount

// logError returns a new error generated by joining 'id' and 'a' and
// optionally calls Sender.Config.LogFunc (if not nil) to log the error.
func (sd *Sender) logError(id uint32, a ...interface{}) error {
	ret := makeError(id, a...)
	if sd.Config != nil && sd.Config.LogFunc != nil {
		msg := ret.Error()
		sd.Config.LogFunc(msg)
	}
	return ret
} //                                                                    logError

// logInfo calls Sender.Config.LogFunc (if not nil) to log a message.
func (sd *Sender) logInfo(a ...interface{}) {
	if sd.Config != nil && sd.Config.LogFunc != nil {
		sd.Config.LogFunc(a...)
	}
} //                                                                     logInfo

// makePacket prepares a packet for immediate sending: it stores,
// hashes data and sets the packet's sentTime to current time.
//
// The size of the packet must not exceed Config.PacketSizeLimit
//
func (sd *Sender) makePacket(data []byte) (*senderPacket, error) {
	if len(data) > sd.Config.PacketSizeLimit {
		return nil, sd.logError(0xE71F9B, "len(data)", len(data),
			"> Config.PacketSizeLimit", sd.Config.PacketSizeLimit)
	}
	sentHash := getHash(data)
	packet := senderPacket{
		data:     data,
		sentHash: sentHash,
		sentTime: time.Now(),
		// confirmedHash, confirmedTime: zero value
	}
	return &packet, nil
} //                                                                  makePacket

// end
