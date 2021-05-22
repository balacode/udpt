// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                          /[const.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// tagDataItemHash tag prefixes a UDP packet sent by the
// sender to request a data item's hash from the receiver.
// This is needed to check if a data item needs sending.
const tagDataItemHash = "HASH:"

// tagFragment prefixes a UDP packet sent by the sender to the receiver,
// containing a fragment of a data item being transferred.
const tagFragment = "FRAG:"

// tagConfirmation tag prefixes a UDP packet sent back by the
// receiver confirming a tagFragment packet sent by the sender.
const tagConfirmation = "CONF:"

// end
