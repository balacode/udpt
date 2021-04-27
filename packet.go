// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                         /[packet.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"time"
)

// Packet _ _
type Packet struct {
	data          []byte
	sentHash      []byte
	sentTime      time.Time
	confirmedHash []byte
	confirmedTime time.Time
} //                                                                      Packet

// NewPacket _ _
func NewPacket(data []byte) (*Packet, error) {
	err := Config.Validate()
	if err != nil {
		return nil, logError(0xEA6BA6, err)
	}
	if len(data) > Config.PacketSizeLimit {
		return nil, logError(0xE71F9B, "len(data)", len(data),
			"> Config.PacketSizeLimit", Config.PacketSizeLimit)
	}
	var (
		tm     = time.Now()
		hash   = getHash(data)
		packet = Packet{data: data, sentHash: hash, sentTime: tm}
		//              confirmedHash & confirmedTime: zero value
	)
	return &packet, nil
} //                                                                   NewPacket

// IsDelivered returns true if a packet has been successfully
// delivered (by receiving a successful confirmation packet).
func (ob *Packet) IsDelivered() bool {
	ret := bytes.Equal(ob.sentHash, ob.confirmedHash)
	return ret
} //                                                                 IsDelivered

// end
