// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                          /[const.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
)

const (
	// DATA_ITEM_HASH tag prefixes a UDP packet sent by the
	// sender to request a data item's hash from the receiver.
	// This is needed to check if a data item needs sending.
	DATA_ITEM_HASH = "HASH:"

	// FRAGMENT tag prefixes a UDP packet sent by the sender to the receiver,
	// containing a fragment of a data item being transferred.
	FRAGMENT = "FRAG:"

	// FRAGMENT_CONFIRMATION tag prefixes a UDP packet sent back by
	// the receiver confirming a FRAGMENT packet sent by the sender.
	FRAGMENT_CONFIRMATION = "CONF:"
)

const (
	// EInvalidArg error indicates an invalid argument.
	EInvalidArg = "Invalid argument"

	// ENilReceiver error indicates a method call on a nil object.
	ENilReceiver = "nil receiver"
)

// pl is fmt.Println() but is used only for debugging.
var pl = fmt.Println

// end
