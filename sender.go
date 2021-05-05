// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                         /[sender.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// # Sender Class
//   Sender struct
//
// # Main Methods (ob *Sender)
//   ) Send(name string, data []byte) error
//   ) SendString(name string, s string) error
//
// # Informatory Properties (ob *Sender)
//   ) AverageResponseMs() float64
//   ) DeliveredAllParts() bool
//   ) TransferSpeedKBpS() float64
//
// # Informatory Methods (ob *Sender)
//   ) PrintInfo()
//
// # Internal Lifecycle Methods (ob *Sender)
//   ) requestDataItemHash(name string) []byte
//   ) connect() error
//   ) sendUndeliveredPackets() error
//   ) collectConfirmations()
//   ) waitForAllConfirmations()
//   ) close() error
//
// # Internal Helper Methods (ob *Sender)
//   ) getPacketCount(length int) int
//   ) makePacket(data []byte) (*Packet, error)
//   ) updateInfo()

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// udpInfo contains UDP transfer statistics, such as the transfer
// speed and the number of packets delivered and lost.
type udpInfo struct {
	averageResponseMs float64
	bytesDelivered    int64
	bytesLost         int64
	packetsDelivered  int64
	packsLost         int64
	transferSpeedKBpS float64
	transferTime      time.Duration
} //                                                                     udpInfo

// udpTotal contains total UDP statistics from time of startup.
// These statistics are accumulated over all Senders after every call to Send.
var udpTotal udpInfo

// -----------------------------------------------------------------------------
// # Sender Class

// Sender is an internal class that coordinates sending
// a sequence of bytes to a listening Receiver.
//
type Sender struct {

	// Address is the domain name or IP address of the
	// listening receiver, excluding the port number.
	Address string

	// Port is the port number of the listening server.
	// This number must be between 1 and 65535.
	Port int

	// CryptoKey is the secret symmetric encryption key that
	// must be shared between the Sender and the Receiver.
	// The correct size of this key depends on
	// the implementation of SymmetricCipher.
	CryptoKey []byte

	// Config contains UDP and other configuration settings.
	// These settings normally don't need to be changed.
	Config ConfigSettings

	// -------------------------------------------------------------------------

	// conn holds the UDP connection to a Receiver
	conn *net.UDPConn

	// dataHash contains the hash of all bytes of the data item being sent
	dataHash []byte

	// info contains UDP transfer statistics, such as the transfer
	// speed and the number of packets delivered and lost
	info udpInfo

	// packets contains all the packets of the currently transferred data item;
	// some of them may have been delivered, while others may need (re)sending
	packets []Packet

	// startTime is the time the first packet was sent, after
	// the bytes of the data item have been compressed
	startTime time.Time

	// wg is used by waitForAllConfirmations()
	// to wait for sendUndeliveredPackets()
	wg sync.WaitGroup
} //                                                                      Sender

// -----------------------------------------------------------------------------
// # Main Methods (ob *Sender)

// Send transfers a sequence of bytes ('data') to the
// Receiver specified by Sender.Address and Port.
func (ob *Sender) Send(name string, data []byte) error {
	//
	err := ob.Config.Validate()
	if err != nil {
		return ob.logError(0xE5D92D, err)
	}
	if strings.TrimSpace(ob.Address) == "" {
		return ob.logError(0xE5A04A, "missing Sender.Address")
	}
	if ob.Port < 1 || ob.Port > 65535 {
		return ob.logError(0xE7B72A, "invalid Sender.Port:", ob.Port)
	}
	if ob.Config.Cipher == nil {
		var aes AESCipher
		err := aes.InitCipher(ob.CryptoKey)
		if err != nil {
			return ob.logError(0xE5EC36, err)
		}
		ob.Config.Cipher = &aes
	}
	err = ob.Config.Cipher.ValidateKey(ob.CryptoKey)
	if err != nil {
		return ob.logError(0xE3A5FF, "invalid Sender.CryptoKey:", err)
	}
	hash, err := getHash(data)
	if err != nil {
		return ob.logError(0xE4B4D8, err)
	}
	if ob.Config.VerboseSender {
		ob.logInfo("\n" + strings.Repeat("-", 80) + "\n" +
			fmt.Sprintf("Send name: %s size: %d hash: %X",
				name, len(data), hash))
	}
	remoteHash := ob.requestDataItemHash(name)
	if bytes.Equal(hash, remoteHash) {
		return nil
	}
	compressed, err := compress(data)
	if err != nil {
		return ob.logError(0xE2A7C3, err)
	}
	packetCount := ob.getPacketCount(len(compressed))
	ob.dataHash, err = getHash(data)
	if err != nil {
		return ob.logError(0xE5E0E6, err)
	}
	ob.startTime = time.Now()
	ob.packets = make([]Packet, packetCount)
	for i := range ob.packets {
		a := i * ob.Config.PacketPayloadSize
		b := a + ob.Config.PacketPayloadSize
		if b > len(compressed) {
			b = len(compressed)
		}
		header := tagFragment + fmt.Sprintf(
			"name:%s hash:%X sn:%d count:%d\n",
			name, ob.dataHash, i+1, packetCount,
		)
		packet, err2 := ob.makePacket(
			append([]byte(header), compressed[a:b]...),
		)
		if err2 != nil {
			return ob.logError(0xE567A4, err2)
		}
		ob.packets[i] = *packet
	}
	newConn, err := ob.connect()
	if err != nil {
		return ob.logError(0xE8B8D0, err)
	}
	ob.conn = newConn
	go ob.collectConfirmations()
	for retries := 0; retries < ob.Config.SendRetries; retries++ {
		err = ob.sendUndeliveredPackets()
		if err != nil {
			defer func() {
				err2 := ob.close()
				if err2 != nil {
					_ = ob.logError(0xE71C7A, err2)
				}
			}()
			return ob.logError(0xE23CE0, err)
		}
		ob.waitForAllConfirmations()
		if ob.DeliveredAllParts() {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
	ob.updateInfo()
	err = ob.close()
	if err != nil {
		return ob.logError(0xE40A05, err)
	}
	if !ob.DeliveredAllParts() {
		return ob.logError(0xE1C3A7, "undelivered packets")
	}
	remoteHash = ob.requestDataItemHash(name)
	if !bytes.Equal(hash, remoteHash) {
		return ob.logError(0xE1F101, "hash mismatch")
	}
	if ob.Config.VerboseSender {
		ob.PrintInfo()
	}
	return nil
} //                                                                        Send

// SendString transfers string 's' to the Receiver
// specified by Sender.Address and Port.
//
func (ob *Sender) SendString(name string, s string) error {
	return ob.Send(name, []byte(s))
} //                                                                  SendString

// -----------------------------------------------------------------------------
// # Informatory Properties (ob *Sender)

// AverageResponseMs is the average response time, in milliseconds, between
// a packet being sent and its delivery confirmation being received.
func (ob *Sender) AverageResponseMs() float64 {
	if ob == nil {
		_ = ob.logError(0xE1B78F, ENilReceiver)
		return 0.0
	}
	if ob.info.packetsDelivered == 0 {
		return 0.0
	}
	// instead of using transferTime.Milliseconds(),
	// cast to float64 to get sub-millisecond timing
	ret := float64(ob.info.transferTime) /
		float64(time.Millisecond) /
		float64(ob.info.packetsDelivered)
	return ret
} //                                                           AverageResponseMs

// DeliveredAllParts returns true if all parts of the
// sent data item have been delivered. I.e. all packets
// have been sent, resent if needed, and confirmed.
func (ob *Sender) DeliveredAllParts() bool {
	if ob == nil {
		_ = ob.logError(0xE52E72, ENilReceiver)
		return false
	}
	ret := true
	for _, packet := range ob.packets {
		if !bytes.Equal(packet.sentHash, packet.confirmedHash) {
			ret = false
			break
		}
	}
	return ret
} //                                                           DeliveredAllParts

// TransferSpeedKBpS returns the transfer speed of the current Send
// operation, in Kilobytes (more accurately, Kibibytes) per second.
func (ob *Sender) TransferSpeedKBpS() float64 {
	if ob == nil {
		_ = ob.logError(0xE6C59B, ENilReceiver)
		return 0.0
	}
	if ob.info.transferTime < 1 {
		return 0.0
	}
	sec := float64(ob.info.transferTime) / float64(time.Second)
	ret := float64(ob.info.bytesDelivered/1024) / sec
	return ret
} //                                                           TransferSpeedKBpS

// -----------------------------------------------------------------------------
// # Informatory Methods (ob *Sender)

// PrintInfo prints the UDP transfer statistics to the standard output.
func (ob *Sender) PrintInfo() {
	if ob == nil {
		_ = ob.logError(0xE483B1, ENilReceiver)
		return
	}
	tItem := time.Duration(0)
	for i, pack := range ob.packets {
		tPack, status := time.Duration(0), "✔"
		if pack.isDelivered() {
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
			ms = padf(9, "%0.1f ms",
				float64(tPack)/float64(time.Millisecond))
		)
		if pack.confirmedTime.IsZero() {
			t1 = "NONE"
		}
		ob.logInfo("SN:", sn, "T0:", t0, "T1:", t1, status, ms)
		tItem += tPack
	}
	var (
		sec          = ob.info.transferTime.Seconds()
		totalSeconds = udpTotal.transferTime.Seconds()
		avg          = ob.AverageResponseMs()
		speed        = ob.TransferSpeedKBpS()
		prt          = func(tag, format string, v1, v2 interface{}) {
			ob.logInfo(tag, padf(12, format, v1), fmt.Sprintf(format, v2))
		}
	)
	prt("B. delivered:", "%d", ob.info.bytesDelivered, udpTotal.bytesDelivered)
	prt("Bytes lost  :", "%d", ob.info.bytesLost, udpTotal.bytesLost)
	prt("P. delivered:", "%d", ob.info.packetsDelivered,
		udpTotal.packetsDelivered)
	prt("Packets lost:", "%d", ob.info.packsLost, udpTotal.packsLost)
	prt("Time in item:", "%0.1f s", sec, totalSeconds)
	prt("Avg./ Packet:", "%0.1f ms", avg, udpTotal.averageResponseMs)
	prt("Trans. speed:", "%0.1f KiB/s", speed, udpTotal.transferSpeedKBpS)
} //                                                                   PrintInfo

// -----------------------------------------------------------------------------
// # Internal Lifecycle Methods (ob *Sender)

// connect connects to the Receiver at Sender.Address and Port
// and returns a new UDP connection and an error value.
//
// Note that it doesn't change the value of Sender.conn
//
func (ob *Sender) connect() (*net.UDPConn, error) {
	if ob == nil {
		return nil, ob.logError(0xE65C26, ENilReceiver)
	}
	addr := fmt.Sprintf("%s:%d", ob.Address, ob.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, ob.logError(0xEC7C6B, err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr) // (*net.UDPConn, error)
	if err != nil {
		return nil, ob.logError(0xE15CE1, err)
	}
	// TODO: add this to ConfigSettings
	err = conn.SetWriteBuffer(16 * 1024 * 2014) // 16 MiB
	if err != nil {
		return nil, ob.logError(0xE5F9C7, err)
	}
	return conn, nil
} //                                                                     connect

// requestDataItemHash requests and waits for the listening receiver to
// return the hash of the data item identified by 'name'. If the receiver
// can locate the data item, it returns its hash, otherwise it returns nil.
func (ob *Sender) requestDataItemHash(name string) []byte {
	err := ob.Config.Validate()
	if err != nil {
		_ = ob.logError(0xE5BC2E, err)
		return nil
	}
	tempConn, err := ob.connect()
	if err != nil {
		_ = ob.logError(0xE7DF8B, err)
		return nil
	}
	packet, err := ob.makePacket([]byte(tagDataItemHash + name))
	if err != nil {
		_ = ob.logError(0xE1F8C5, err)
		return nil
	}
	err = packet.send(tempConn, ob.Config.Cipher)
	if err != nil {
		_ = ob.logError(0xE7F316, err)
		return nil
	}
	encryptedReply := make([]byte, ob.Config.PacketSizeLimit)
	nRead, _, err :=
		readFromUDPConn(tempConn, encryptedReply, ob.Config.ReplyTimeout)
	if err != nil {
		_ = ob.logError(0xE97FC3, err)
		return nil
	}
	reply, err := ob.Config.Cipher.Decrypt(encryptedReply[:nRead])
	if err != nil {
		_ = ob.logError(0xE2B5A1, err)
		return nil
	}
	var hash []byte
	if len(reply) > 0 {
		if !bytes.HasPrefix(reply, []byte(tagDataItemHash)) {
			_ = ob.logError(0xE08AD4, "invalid reply:", reply)
			return nil
		}
		hexHash := string(reply[len(tagDataItemHash):])
		if hexHash == "not_found" {
			return nil
		}
		hash, err = hex.DecodeString(hexHash)
		if err != nil {
			_ = ob.logError(0xE5A4E7, err)
			return nil
		}
	}
	return hash
} //                                                         requestDataItemHash

// sendUndeliveredPackets sends all undelivered
// packets to the destination Receiver.
func (ob *Sender) sendUndeliveredPackets() error {
	if ob == nil {
		return ob.logError(0xE8DB3F, ENilReceiver)
	}
	err := ob.Config.Validate()
	if err != nil {
		return ob.logError(0xE86B5B, err)
	}
	n := len(ob.packets)
	for i := 0; i < n; i++ {
		packet := &ob.packets[i]
		if packet.isDelivered() {
			continue
		}
		time.Sleep(2 * time.Millisecond)
		ob.wg.Add(1)
		go func() {
			err := packet.send(ob.conn, ob.Config.Cipher)
			if err != nil {
				_ = ob.logError(0xE67BA4, err)
			}
			ob.wg.Done()
		}()
	}
	return nil
} //                                                      sendUndeliveredPackets

// collectConfirmations enters a loop that receives confirmation packets
// from the sender, and marks all confirmed packets as delivered.
func (ob *Sender) collectConfirmations() {
	if ob == nil {
		_ = ob.logError(0xE8EA91, ENilReceiver)
		return
	}
	err := ob.Config.Validate()
	if err != nil {
		_ = ob.logError(0xE44C4A, err)
		return
	}
	encryptedReply := make([]byte, ob.Config.PacketSizeLimit)
	for ob.conn != nil {
		// 'encryptedReply' is overwritten after every readFromUDPConn
		nRead, addr, err :=
			readFromUDPConn(ob.conn, encryptedReply, ob.Config.ReplyTimeout)
		if err == errClosed {
			break
		}
		if err != nil {
			_ = ob.logError(0xE7B6B2, err)
			continue
		}
		if nRead == 0 {
			_ = ob.logError(0xE4CB0B, "received no data")
			continue
		}
		recv, err := ob.Config.Cipher.Decrypt(encryptedReply[:nRead])
		if err != nil {
			_ = ob.logError(0xE5C43E, err)
			continue
		}
		if !bytes.HasPrefix(recv, []byte(tagConfirmation)) {
			_ = ob.logError(0xE5AF24, "bad reply header")
			if ob.Config.VerboseSender {
				ob.logInfo("ERROR received:", len(recv), "bytes")
			}
			continue
		}
		confirmedHash := recv[len(tagConfirmation):]
		if ob.Config.VerboseSender {
			ob.logInfo("Sender received", nRead, "bytes from", addr)
		}
		go func(confirmedHash []byte) {
			for i, packet := range ob.packets {
				if bytes.Equal(packet.sentHash, confirmedHash) {
					ob.packets[i].confirmedTime = time.Now()
					ob.packets[i].confirmedHash = confirmedHash
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
func (ob *Sender) waitForAllConfirmations() {
	if ob == nil {
		_ = ob.logError(0xE2A34E, ENilReceiver)
		return
	}
	err := ob.Config.Validate()
	if err != nil {
		_ = ob.logError(0xE4B72B, err)
		return
	}
	ob.logInfo("Waiting . . .")
	t0 := time.Now()
	ob.wg.Wait()
	for {
		time.Sleep(50 * time.Millisecond)
		if ob.DeliveredAllParts() {
			if ob.Config.VerboseSender {
				ob.logInfo("Delivered all packets")
			}
			break
		}
		since := time.Since(t0)
		if since >= ob.Config.ReplyTimeout {
			ob.logInfo("Config.ReplyTimeout exceeded",
				fmt.Sprintf("%0.1f", since.Seconds()))
			break
		}
	}
	for _, packet := range ob.packets {
		if packet.isDelivered() {
			ob.info.bytesDelivered += int64(len(packet.data))
			ob.info.packetsDelivered++
		} else {
			ob.info.bytesLost += int64(len(packet.data))
			ob.info.packsLost++
		}
	}
	if ob.Config.VerboseSender {
		ob.logInfo("Waited:", time.Since(t0))
	}
} //                                                     waitForAllConfirmations

// close closes the UDP connection.
func (ob *Sender) close() error {
	if ob == nil {
		return ob.logError(0xE0561D, ENilReceiver)
	}
	err := ob.conn.Close()
	ob.conn = nil
	if err != nil {
		return ob.logError(0xE71AB2, err)
	}
	return nil
} //                                                                       close

// -----------------------------------------------------------------------------
// # Internal Helper Methods (ob *Sender)

// getPacketCount calculates the number of packets needed to send 'length'
// bytes. This depends on the setting of Config.PacketPayloadSize.
func (ob *Sender) getPacketCount(length int) int {
	err := ob.Config.Validate()
	if err != nil {
		_ = ob.logError(0xEC866E, err)
		return 0
	}
	if length < 1 {
		return 0
	}
	count := length / ob.Config.PacketPayloadSize
	if (count * ob.Config.PacketPayloadSize) < length {
		count++
	}
	return count
} //                                                              getPacketCount

// logError returns a new error value generated by joining id and args
// and optionally calls Sender.LogFunc (if not nil) to log the error.
func (ob *Sender) logError(id uint32, args ...interface{}) error {
	ret := makeError(id, args...)
	if ob.Config.LogFunc != nil {
		msg := ret.Error()
		ob.Config.LogFunc(msg)
	}
	return ret
} //                                                                    logError

// logInfo calls Sender.LogFunc (if not nil) to log a message.
func (ob *Sender) logInfo(args ...interface{}) {
	if ob.Config.LogFunc != nil {
		ob.Config.LogFunc(args...)
	}
} //                                                                     logInfo

// makePacket prepares a packet for immediate sending: it stores,
// hashes data and sets the packet's sentTime to current time.
//
// The size of the packet must not exceed Config.PacketSizeLimit
//
func (ob *Sender) makePacket(data []byte) (*Packet, error) {
	if len(data) > ob.Config.PacketSizeLimit {
		return nil, ob.logError(0xE71F9B, "len(data)", len(data),
			"> Config.PacketSizeLimit", ob.Config.PacketSizeLimit)
	}
	sentHash, err := getHash(data)
	if err != nil {
		return nil, ob.logError(0xE84C0B, err)
	}
	packet := Packet{
		data:     data,
		sentHash: sentHash,
		sentTime: time.Now(),
		// confirmedHash, confirmedTime: zero value
	}
	return &packet, nil
} //                                                                  makePacket

// updateInfo updates the global UDP transfer statistics
// with the statistics of the current Send operation.
func (ob *Sender) updateInfo() {
	if ob == nil {
		_ = ob.logError(0xED48D1, ENilReceiver)
		return
	}
	ob.info.transferTime = time.Since(ob.startTime)
	//
	// update global statistics
	udpTotal.bytesDelivered += ob.info.bytesDelivered
	udpTotal.bytesLost += ob.info.bytesLost
	udpTotal.packetsDelivered += ob.info.packetsDelivered
	udpTotal.packsLost += ob.info.packsLost
	udpTotal.transferTime += ob.info.transferTime
	n := 0.0
	if udpTotal.packetsDelivered > 0 {
		ms := float64(udpTotal.transferTime) / float64(time.Millisecond)
		n = ms / float64(udpTotal.packetsDelivered)
	}
	udpTotal.averageResponseMs = n
	n = 0.0
	if udpTotal.transferTime > 0 {
		sec := float64(udpTotal.transferTime) / float64(time.Second)
		n = float64(udpTotal.bytesDelivered/1024) / sec
	}
	udpTotal.transferSpeedKBpS = n
} //                                                                  updateInfo

// end
