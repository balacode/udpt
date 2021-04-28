// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                      /[data_item.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
)

// DataItem holds a data item being received by the Receiver. A data item
// is just a sequence of bytes being transferred. It could be a file,
// a JSON string or any other resource.
//
// Since we're using UDP, which has a limited packet size, the resource
// is split into several smaller pieces that are sent as UDP packets.
//
type DataItem struct {
	Name                 string
	Hash                 []byte
	CompressedPieces     [][]byte
	CompressedSizeInfo   int
	UncompressedSizeInfo int
} //                                                                    DataItem

// -----------------------------------------------------------------------------
// # Property

// IsLoaded returns true if the current data item has been
// fully received (all its pieces have been collected).
func (ob *DataItem) IsLoaded() bool {
	ret := true
	for _, piece := range ob.CompressedPieces {
		if len(piece) < 1 {
			ret = false
			break
		}
	}
	return ret
} //                                                                    IsLoaded

// -----------------------------------------------------------------------------
// # Methods

// PrintInfo prints information on the current data item
func (ob *DataItem) PrintInfo(tag string) {
	logInfo(tag+" name:", ob.Name)
	logInfo(tag+" hash:", ob.Hash)
	logInfo(tag+" pcs.:", len(ob.CompressedPieces))
	logInfo(tag+" comp:", ob.CompressedSizeInfo, "bytes")
	logInfo(tag+" size:", ob.UncompressedSizeInfo, "bytes")
} //                                                                   PrintInfo

// Reset discards the contents of the data item and clears its name and hash.
func (ob *DataItem) Reset() {
	ob.Name = ""
	ob.Hash = nil
	ob.CompressedPieces = nil
	ob.CompressedSizeInfo = 0
	ob.UncompressedSizeInfo = 0
} //                                                                       Reset

// Retain changes the Name, Hash and empties CompressedPieces if the passed
// name, hash, and packetCount don't match their existing values in the object.
func (ob *DataItem) Retain(name string, hash []byte, packetCount int) {
	if ob.Name == name &&
		bytes.Equal(ob.Hash, hash) &&
		len(ob.CompressedPieces) == packetCount {
		return
	}
	ob.Name = name
	ob.Hash = hash
	ob.CompressedPieces = make([][]byte, packetCount)
	ob.CompressedSizeInfo = 0
	ob.UncompressedSizeInfo = 0
} //                                                                      Retain

// UnpackBytes _ _
func (ob *DataItem) UnpackBytes() ([]byte, error) {
	//
	// join pieces (provided all have been collected) to get compressed data
	if !ob.IsLoaded() {
		return nil, logError(0xE76AF5, ": data item is incomplete")
	}
	compressed := bytes.Join(ob.CompressedPieces, nil)
	ob.CompressedSizeInfo = len(compressed)
	//
	// uncompress data
	ret, err := uncompress(compressed)
	if err != nil {
		return nil, logError(0xE95DFB, "(uncompress):", err)
	}
	ob.UncompressedSizeInfo = len(ret)
	//
	// hash of uncompressed data should match original hash
	hash := getHash(ret)
	if !bytes.Equal(hash, ob.Hash) {
		return nil, logError(0xE87D89, ": checksum mismatch")
	}
	return ret, nil
} //                                                                 UnpackBytes

// end
