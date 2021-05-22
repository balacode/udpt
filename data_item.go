// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                      /[data_item.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"fmt"
	"io"
)

// dataItem holds a data item being received by a Receiver. A data item
// is just a sequence of bytes being transferred. It could be a file,
// a JSON string or any other resource.
//
// Since we're using UDP, which has a limited packet size, the resource
// is split into several smaller pieces that are sent as UDP packets.
//
type dataItem struct {
	Name                 string
	Hash                 []byte
	CompressedPieces     [][]byte
	CompressedSizeInfo   int
	UncompressedSizeInfo int
} //                                                                    dataItem

// -----------------------------------------------------------------------------
// # Property

// IsLoaded returns true if the current data item has been
// fully received (all its pieces have been collected).
//
// If the item has no pieces, returns false.
//
func (di *dataItem) IsLoaded() bool {
	ret := len(di.CompressedPieces) > 0
	for _, piece := range di.CompressedPieces {
		if len(piece) < 1 {
			ret = false
			break
		}
	}
	return ret
} //                                                                    IsLoaded

// -----------------------------------------------------------------------------
// # Methods

// LogStats writes details of the current data item to the
// passed io.Writer. Each written line is prefixed with tag.
//
func (di *dataItem) LogStats(tag string, w io.Writer) {
	log := func(v ...interface{}) {
		s := fmt.Sprintln(v...)
		w.Write([]byte(s))
	}
	log(tag, "name:", di.Name)
	log(tag, "hash:", fmt.Sprintf("%X", di.Hash))
	log(tag, "pcs.:", len(di.CompressedPieces))
	log(tag, "comp:", di.CompressedSizeInfo, "bytes")
	log(tag, "size:", di.UncompressedSizeInfo, "bytes")
} //                                                                    LogStats

// Reset discards the contents of the data item and clears its name and hash.
func (di *dataItem) Reset() {
	di.Name = ""
	di.Hash = nil
	di.CompressedPieces = nil
	di.CompressedSizeInfo = 0
	di.UncompressedSizeInfo = 0
} //                                                                       Reset

// Retain changes the Name, Hash, and empties CompressedPieces when the passed
// name, hash and packetCount don't match their current values in the object.
func (di *dataItem) Retain(name string, hash []byte, packetCount int) {
	if di.Name == name &&
		bytes.Equal(di.Hash, hash) &&
		len(di.CompressedPieces) == packetCount {
		return
	}
	di.Name = name
	di.Hash = hash
	di.CompressedPieces = make([][]byte, packetCount)
	di.CompressedSizeInfo = 0
	di.UncompressedSizeInfo = 0
} //                                                                      Retain

// UnpackBytes joins CompressedPieces and uncompresses
// the resulting bytes to get the original data item.
func (di *dataItem) UnpackBytes(compressor Compression) ([]byte, error) {
	//
	// join pieces (provided all have been collected) to get compressed data
	if !di.IsLoaded() {
		return nil, makeError(0xE76AF5, "data item is incomplete")
	}
	compressed := bytes.Join(di.CompressedPieces, nil)
	di.CompressedSizeInfo = len(compressed)
	//
	// uncompress data
	ret, err := compressor.Uncompress(compressed)
	if err != nil {
		return nil, makeError(0xE95DFB, err)
	}
	di.UncompressedSizeInfo = len(ret)
	//
	// hash of uncompressed data should match original hash
	hash := getHash(ret)
	if !bytes.Equal(hash, di.Hash) {
		return nil, makeError(0xE87D89, "hash mismatch")
	}
	return ret, nil
} //                                                                 UnpackBytes

// end
