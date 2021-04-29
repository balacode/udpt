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

// IsDelivered returns true if a packet has been successfully
// delivered (by receiving a successful confirmation packet).
func (ob *Packet) IsDelivered() bool {
	ret := bytes.Equal(ob.sentHash, ob.confirmedHash)
	return ret
} //                                                                 IsDelivered

// end
