// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                         /[sender.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// # Sender Class
//   Sender struct
//
// # Public Methods
//   ) Send(name string, data []byte) error
//   ) SendString(name string, s string) error
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
//
// # Information Properties
//   ) averageResponseMs() float64
//   ) deliveredAllParts() bool
//   ) transferSpeedKBpS() float64
//
// # Information Methods
//   ) printInfo()

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// udpInfo contains global UDP transfer statistics since startup.
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
// These statistics are accumulated after every call to Send.
var udpTotal udpInfo

// -----------------------------------------------------------------------------
// # Sender Class

// Sender is an internal class that coordinates sending
// a sequence of bytes to a listening Receiver.
//
type Sender struct {

	// Address is the domain name or IP address of the listening receiver,
	// excluding the port number.
	Address string

	// Port is the port number of the listening server.
	// This number must be between 1 and 65535.
	Port int

	// CryptoKey is the secret symmetric encryption key that
	// must be shared between the sender and the receiver.
	// The correct size of this key depends
	// on the implementation of SymmetricCipher.
	CryptoKey []byte

	// Config _ _
	Config ConfigSettings

	// dataHash _ _
	dataHash []byte

	// startTime _ _
	startTime time.Time

	// packets _ _
	packets []Packet

	// conn _ _
	conn *net.UDPConn

	// wg _ _
	wg sync.WaitGroup

	// info _ _
	info udpInfo
} //                                                                      Sender

// -----------------------------------------------------------------------------
// # Public Methods

// Send sends (transfers) a sequence of bytes ('data') to the
// Receiver specified by Sender.Address and Sender.Port.
//
func (ob *Sender) Send(name string, data []byte) error {
	//
	if strings.TrimSpace(ob.Address) == "" {
		return logError(0xE5A04A, "missing Address")
	}
	if ob.Port < 1 || ob.Port > 65535 {
		return logError(0xE7B72A, "invalid Port:", ob.Port)
	}
	if len(ob.CryptoKey) != 32 {
		return logError(0xEB8484, "CryptoKey must be 32, but is",
			len(ob.CryptoKey), "bytes long")
	}
	err := ob.Config.Validate()
	if err != nil {
		return logError(0xE5D92D, err)
	}
	hash := getHash(data)
	if ob.Config.VerboseSender {
		logInfo("\n" + strings.Repeat("-", 80) + "\n" +
			fmt.Sprintf("Send name: %s size: %d hash: %X",
				name, len(data), hash))
	}
	remoteHash := ob.requestDataItemHash(name)
	if bytes.Equal(hash, remoteHash) {
		return nil
	}
	compressed, err := compress(data)
	if err != nil {
		return logError(0xE2A7C3, "(compress):", err)
	}
	packetCount := ob.getPacketCount(len(compressed))
	ob.dataHash = getHash(data)
	ob.startTime = time.Now()
	ob.packets = make([]Packet, packetCount)
	for i := range ob.packets {
		a := i * ob.Config.PacketPayloadSize
		b := a + ob.Config.PacketPayloadSize
		if b > len(compressed) {
			b = len(compressed)
		}
		header := FRAGMENT + fmt.Sprintf(
			"name:%s hash:%X sn:%d count:%d\n",
			name, ob.dataHash, i+1, packetCount,
		)
		packet, err2 := ob.makePacket(
			append([]byte(header), compressed[a:b]...),
		)
		if err2 != nil {
			return logError(0xE567A4, "(makePacket):", err2)
		}
		ob.packets[i] = *packet
	}
	err = ob.connect()
	if err != nil {
		return logError(0xE8B8D0, "(connect):", err)
	}
	go ob.collectConfirmations()
	for retries := 0; retries < ob.Config.SendRetries; retries++ {
		err = ob.sendUndeliveredPackets()
		if err != nil {
			defer func() {
				err2 := ob.close()
				if err2 != nil {
					_ = logError(0xE71C7A, "(close):", err2)
				}
			}()
			return logError(0xE23CE0, "(sendUndeliveredPackets):", err)
		}
		ob.waitForAllConfirmations()
		if ob.deliveredAllParts() {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
	ob.updateInfo()
	err = ob.close()
	if err != nil {
		return logError(0xE40A05, "(close):", err)
	}
	if !ob.deliveredAllParts() {
		return logError(0xE1C3A7, ": undelivered packets")
	}
	remoteHash = ob.requestDataItemHash(name)
	if !bytes.Equal(hash, remoteHash) {
		return logError(0xE1F101, ": hash mismatch")
	}
	ob.printInfo()
	return nil
} //                                                                        Send

// SendString sends (transfers) string 's' to the Receiver
// specified by Sender.Address and Sender.Port.
//
func (ob *Sender) SendString(name string, s string) error {
	return ob.Send(name, []byte(s))
} //                                                                  SendString

// -----------------------------------------------------------------------------
// # Internal Lifecycle Methods (ob *Sender)

// requestDataItemHash requests and waits for the listening receiver
// to return the hash of the data item named by 'name'. If the receiver
// can locate the data item, returns its hash, otherwise returns nil.
func (ob *Sender) requestDataItemHash(name string) []byte {
	err := ob.Config.Validate()
	if err != nil {
		_ = logError(0xE5BC2E, err)
		return nil
	}
	addr := fmt.Sprintf("%s:%d", ob.Address, ob.Port)
	conn, err := connect(addr)
	if err != nil {
		_ = logError(0xE7DF8B, "(connect):", err)
		return nil
	}
	packet, err := ob.makePacket([]byte(DATA_ITEM_HASH + name))
	if err != nil {
		_ = logError(0xE1F8C5, "(makePacket):", err)
		return nil
	}
	err = sendPacket(packet, ob.CryptoKey, conn)
	if err != nil {
		_ = logError(0xE7F316, "(sendPacket):", err)
		return nil
	}
	encryptedReply := make([]byte, ob.Config.PacketSizeLimit)
	nRead, _ /*addr*/, err :=
		readFromUDPConn(conn, encryptedReply, ob.Config.ReplyTimeout)
	if err != nil {
		_ = logError(0xE97FC3, "(ReadFrom):", err)
		return nil
	}
	reply, err := aesDecrypt(encryptedReply[:nRead], ob.CryptoKey)
	if err != nil {
		_ = logError(0xE2B5A1, "(aesDecrypt):", err)
		return nil
	}
	var hash []byte
	if len(reply) > 0 {
		if !bytes.HasPrefix(reply, []byte(DATA_ITEM_HASH)) {
			_ = logError(0xE08AD4, ": invalid reply:", reply)
			return nil
		}
		hexHash := string(reply[len(DATA_ITEM_HASH):])
		if hexHash == "not_found" {
			return nil
		}
		hash, err = hex.DecodeString(hexHash)
		if err != nil {
			_ = logError(0xE5A4E7, "(hex.DecodeString):", err)
			return nil
		}
	}
	return hash
} //                                                         requestDataItemHash

// connect connects to the Receiver at Sender.Address and Sender.Port
func (ob *Sender) connect() error {
	if ob == nil {
		return logError(0xE65C26, ":", ENilReceiver)
	}
	addr := fmt.Sprintf("%s:%d", ob.Address, ob.Port)
	conn, err := connect(addr)
	if err != nil {
		ob.conn = nil
		return logError(0xE95B4D, "(connect):", err)
	}
	ob.conn = conn
	return nil
} //                                                                     connect

// sendUndeliveredPackets sends all undelivered
// packets to the destination Receiver.
func (ob *Sender) sendUndeliveredPackets() error {
	if ob == nil {
		return logError(0xE8DB3F, ":", ENilReceiver)
	}
	err := ob.Config.Validate()
	if err != nil {
		return logError(0xE86B5B, err)
	}
	n := len(ob.packets)
	for i := 0; i < n; i++ {
		packet := &ob.packets[i]
		if packet.IsDelivered() {
			continue
		}
		time.Sleep(2 * time.Millisecond)
		ob.wg.Add(1)
		go func() {
			err := sendPacket(packet, ob.CryptoKey, ob.conn)
			if err != nil {
				_ = logError(0xE67BA4, "(sendPacket):", err)
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
		_ = logError(0xE8EA91, ":", ENilReceiver)
		return
	}
	err := ob.Config.Validate()
	if err != nil {
		_ = logError(0xE44C4A, err)
		return
	}
	encryptedReply := make([]byte, ob.Config.PacketSizeLimit)
	for {
		// 'encryptedReply' is overwritten after every readFromUDPConn
		nRead, addr, err :=
			readFromUDPConn(ob.conn, encryptedReply, ob.Config.ReplyTimeout)
		if err != nil {
			if strings.Contains(err.Error(), "closed network connection") {
				return
			}
			_ = logError(0xE7B6B2, "(ReadFrom):", err)
			continue
		}
		if nRead == 0 {
			_ = logError(0xE4CB0B, ": received no data")
			continue
		}
		recv, err := aesDecrypt(encryptedReply[:nRead], ob.CryptoKey)
		if err != nil {
			_ = logError(0xE5C43E, "(aesDecrypt):", err)
			continue
		}
		if !bytes.HasPrefix(recv, []byte(FRAGMENT_CONFIRMATION)) {
			_ = logError(0xE5AF24, ": bad reply header")
			if ob.Config.VerboseSender {
				logInfo("ERROR received:", len(recv), "bytes")
			}
			continue
		}
		confirmedHash := recv[len(FRAGMENT_CONFIRMATION):]
		if ob.Config.VerboseSender {
			logInfo("Sender received", nRead, "bytes from", addr)
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
// will only wait for the duration specified in Config.ReplyTimeout
//
func (ob *Sender) waitForAllConfirmations() {
	if ob == nil {
		_ = logError(0xE2A34E, ":", ENilReceiver)
		return
	}
	err := ob.Config.Validate()
	if err != nil {
		_ = logError(0xE4B72B, err)
		return
	}
	logInfo("Waiting . . .")
	t0 := time.Now()
	ob.wg.Wait()
	for {
		time.Sleep(50 * time.Millisecond)
		ok := ob.deliveredAllParts()
		if ok {
			if ob.Config.VerboseSender {
				logInfo("Delivered all packets")
			}
			break
		}
		since := time.Since(t0)
		if since >= ob.Config.ReplyTimeout {
			logInfo("Config.ReplyTimeout exceeded",
				fmt.Sprintf("%0.1f", since.Seconds()))
			break
		}
	}
	for _, packet := range ob.packets {
		if packet.IsDelivered() {
			ob.info.bytesDelivered += int64(len(packet.data))
			ob.info.packetsDelivered++
		} else {
			ob.info.bytesLost += int64(len(packet.data))
			ob.info.packsLost++
		}
	}
	if ob.Config.VerboseSender {
		logInfo("Waited:", time.Since(t0))
	}
} //                                                     waitForAllConfirmations

// close closes the UDP connection.
func (ob *Sender) close() error {
	if ob == nil {
		return logError(0xE0561D, ":", ENilReceiver)
	}
	err := ob.conn.Close()
	if err != nil {
		ob.conn = nil
		return logError(0xE71AB2, "(close):", err)
	}
	return nil
} //                                                                       close

// -----------------------------------------------------------------------------
// # Internal Helper Methods (ob *Sender)

// getPacketCount calculates the number of packets needed to send 'length'
// bytes. This depends on the setting of Config.PacketPayloadSize.
//
func (ob *Sender) getPacketCount(length int) int {
	err := ob.Config.Validate()
	if err != nil {
		_ = logError(0xEC866E, err)
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

// makePacket _ _
func (ob *Sender) makePacket(data []byte) (*Packet, error) {
	if len(data) > ob.Config.PacketSizeLimit {
		return nil, logError(0xE71F9B, "len(data)", len(data),
			"> Config.PacketSizeLimit", ob.Config.PacketSizeLimit)
	}
	packet := Packet{
		data:     data,
		sentHash: getHash(data),
		sentTime: time.Now(),
		// confirmedHash, confirmedTime: zero value
	}
	return &packet, nil
} //                                                                  makePacket

// -----------------------------------------------------------------------------
// # Information Properties

// averageResponseMs is the average response time, in milliseconds,
// between a packet being sent and a confirmation being received.
func (ob *Sender) averageResponseMs() float64 {
	if ob == nil {
		_ = logError(0xE1B78F, ":", ENilReceiver)
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
} //                                                           averageResponseMs

// deliveredAllParts returns true if all parts of the sent
// data item have been delivered. I.e. all packets
// have been sent, resent if needed, and confirmed.
//
func (ob *Sender) deliveredAllParts() bool {
	if ob == nil {
		_ = logError(0xE52E72, ":", ENilReceiver)
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
} //                                                           deliveredAllParts

// transferSpeedKBpS returns the transfer speed of the current Send
// operation, in Kilobytes (more accurately, Kibibytes) per Second.
func (ob *Sender) transferSpeedKBpS() float64 {
	if ob == nil {
		_ = logError(0xE6C59B, ":", ENilReceiver)
		return 0.0
	}
	if ob.info.transferTime < 1 {
		return 0.0
	}
	sec := float64(ob.info.transferTime) / float64(time.Second)
	ret := float64(ob.info.bytesDelivered/1024) / sec
	return ret
} //                                                           transferSpeedKBpS

// -----------------------------------------------------------------------------
// # Information Methods

// printInfo prints the UDP transfer statistics to the standard output.
func (ob *Sender) printInfo() {
	if ob == nil {
		_ = logError(0xE483B1, ":", ENilReceiver)
		return
	}
	tItem := time.Duration(0)
	for i, pack := range ob.packets {
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
			ms = padf(9, "%0.1f ms",
				float64(tPack)/float64(time.Millisecond))
		)
		if pack.confirmedTime.IsZero() {
			t1 = "NONE"
		}
		logInfo("SN:", sn, "T0:", t0, "T1:", t1, status, ms)
		tItem += tPack
	}
	var (
		sec          = ob.info.transferTime.Seconds()
		totalSeconds = udpTotal.transferTime.Seconds()
		avg          = ob.averageResponseMs()
		speed        = ob.transferSpeedKBpS()
		prt          = func(tag, format string, v1, v2 interface{}) {
			logInfo(tag, padf(12, format, v1), fmt.Sprintf(format, v2))
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
} //                                                                   printInfo

// updateInfo updates the global UDP transfer statistics
// with the statistics of the current Send operation.
func (ob *Sender) updateInfo() {
	if ob == nil {
		_ = logError(0xED48D1, ":", ENilReceiver)
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
