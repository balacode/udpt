// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                          /[const.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// tagFragment prefixes a UDP packet sent by the sender to the receiver,
// containing a fragment of a data item being transferred.
const tagFragment = "FRAG:"

// tagConfirmation tag prefixes a UDP packet sent back by the
// receiver confirming a tagFragment packet sent by the sender.
const tagConfirmation = "CONF:"

// end
