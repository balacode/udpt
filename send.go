// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                           /[send.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// Send(name string, data []byte) error
//
// # udpSender Helper Type
//   udpSender struct
//
// # Methods (ob *udpSender)
//   ) connect() error
//   ) sendUndeliveredPackets() error
//   ) collectConfirmations()
//   ) waitForAllConfirmations()
//   ) close() error
//
// # Information Properties
//   ) averageResponseMs() float64
//   ) deliveredAllParts() bool
//   ) transferSpeedKBpS() float64
//
// # Information Methods
//   ) printInfo()
//
// # Functions
//   getPacketCount(length int) int

import (
	"bytes"
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

// Send sends (transfers) a sequence of bytes ('data') to the
// Receiver specified by Config.Address and Config.Port.
//
func Send(name string, data []byte) error {
	err := Config.Validate()
	if err != nil {
		return logError(0xE5D92D, err)
	}
	hash := getHash(data)
	if Config.VerboseSender {
		logInfo("\n" + strings.Repeat("-", 80) + "\n" +
			fmt.Sprintf("Send name: %s size: %d hash: %X",
				name, len(data), hash))
	}
	remoteHash := requestDataItemHash(name)
	if bytes.Equal(hash, remoteHash) {
		return nil
	}
	compressed, err := compress(data)
	if err != nil {
		return logError(0xE2A7C3, "(compress):", err)
	}
	packetCount := getPacketCount(len(compressed))
	sender := udpSender{
		dataHash:  getHash(data),
		startTime: time.Now(),
		packets:   make([]Packet, packetCount),
	}
	for i := range sender.packets {
		a := i * Config.PacketPayloadSize
		b := a + Config.PacketPayloadSize
		if b > len(compressed) {
			b = len(compressed)
		}
		header := FRAGMENT + fmt.Sprintf(
			"name:%s hash:%X sn:%d count:%d\n",
			name, sender.dataHash, i+1, packetCount,
		)
		packet, err2 := NewPacket(append([]byte(header), compressed[a:b]...))
		if err2 != nil {
			return logError(0xE567A4, "(NewPacket):", err2)
		}
		sender.packets[i] = *packet
	}
	err = sender.connect()
	if err != nil {
		return logError(0xE8B8D0, "(connect):", err)
	}
	go sender.collectConfirmations()
	for retries := 0; retries < Config.SendRetries; retries++ {
		err = sender.sendUndeliveredPackets()
		if err != nil {
			defer func() {
				err2 := sender.close()
				if err2 != nil {
					_ = logError(0xE71C7A, "(close):", err2)
				}
			}()
			return logError(0xE23CE0, "(sendUndeliveredPackets):", err)
		}
		sender.waitForAllConfirmations()
		if sender.deliveredAllParts() {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
	sender.updateInfo()
	err = sender.close()
	if err != nil {
		return logError(0xE40A05, "(close):", err)
	}
	if !sender.deliveredAllParts() {
		return logError(0xE1C3A7, ": undelivered packets")
	}
	remoteHash = requestDataItemHash(name)
	if !bytes.Equal(hash, remoteHash) {
		return logError(0xE1F101, ": hash mismatch")
	}
	sender.printInfo()
	return nil
} //                                                                        Send

// -----------------------------------------------------------------------------
// # udpSender Helper Type

// udpSender is an internal class that coordinates sending
// a sequence of bytes to a listening Receiver.
//
type udpSender struct {
	dataHash  []byte
	startTime time.Time
	packets   []Packet
	conn      *net.UDPConn
	wg        sync.WaitGroup
	info      udpInfo
} //                                                                   udpSender

// -----------------------------------------------------------------------------
// # Methods (ob *udpSender)

// connect connects to the Receiver specified
// by Config.Address and Config.Port
//
func (ob *udpSender) connect() error {
	if ob == nil {
		return logError(0xE65C26, ":", ENilReceiver)
	}
	conn, err := connect()
	if err != nil {
		ob.conn = nil
		return logError(0xE95B4D, "(connect):", err)
	}
	ob.conn = conn
	return nil
} //                                                                     connect

// sendUndeliveredPackets sends all undelivered
// packets to the destination Receiver.
func (ob *udpSender) sendUndeliveredPackets() error {
	if ob == nil {
		return logError(0xE8DB3F, ":", ENilReceiver)
	}
	err := Config.Validate()
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
			err := sendPacket(packet, ob.conn)
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
func (ob *udpSender) collectConfirmations() {
	if ob == nil {
		_ = logError(0xE8EA91, ":", ENilReceiver)
		return
	}
	err := Config.Validate()
	if err != nil {
		_ = logError(0xE44C4A, err)
		return
	}
	encryptedReply := make([]byte, Config.PacketSizeLimit)
	for {
		// 'encryptedReply' is overwritten after every readFromUDPConn
		nRead, addr, err := readFromUDPConn(ob.conn, encryptedReply)
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
		recv, err := aesDecrypt(encryptedReply[:nRead], Config.AESKey)
		if err != nil {
			_ = logError(0xE5C43E, "(aesDecrypt):", err)
			continue
		}
		if !bytes.HasPrefix(recv, []byte(FRAGMENT_CONFIRMATION)) {
			_ = logError(0xE5AF24, ": bad reply header")
			if Config.VerboseSender {
				logInfo("ERROR received:", len(recv), "bytes")
			}
			continue
		}
		confirmedHash := recv[len(FRAGMENT_CONFIRMATION):]
		if Config.VerboseSender {
			logInfo("udpSender received", nRead, "bytes from", addr)
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
func (ob *udpSender) waitForAllConfirmations() {
	if ob == nil {
		_ = logError(0xE2A34E, ":", ENilReceiver)
		return
	}
	err := Config.Validate()
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
			if Config.VerboseSender {
				logInfo("Delivered all packets")
			}
			break
		}
		since := time.Since(t0)
		if since >= Config.ReplyTimeout {
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
	if Config.VerboseSender {
		logInfo("Waited:", time.Since(t0))
	}
} //                                                     waitForAllConfirmations

// close closes the UDP connection.
func (ob *udpSender) close() error {
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
// # Information Properties

// averageResponseMs is the average response time, in milliseconds,
// between a packet being sent and a confirmation being received.
func (ob *udpSender) averageResponseMs() float64 {
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
func (ob *udpSender) deliveredAllParts() bool {
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
func (ob *udpSender) transferSpeedKBpS() float64 {
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
func (ob *udpSender) printInfo() {
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
func (ob *udpSender) updateInfo() {
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

// -----------------------------------------------------------------------------
// # Functions

// getPacketCount calculates the number of packets needed to send 'length'
// bytes. This depends on the setting of Config.PacketPayloadSize.
//
func getPacketCount(length int) int {
	err := Config.Validate()
	if err != nil {
		_ = logError(0xEC866E, err)
		return 0
	}
	if length < 1 {
		return 0
	}
	count := length / Config.PacketPayloadSize
	if (count * Config.PacketPayloadSize) < length {
		count++
	}
	return count
} //                                                              getPacketCount

// end
