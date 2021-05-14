// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                      /[data_item.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"fmt"
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
func (ob *dataItem) IsLoaded() bool {
	ret := len(ob.CompressedPieces) > 0
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

// LogStats prints details of the current data item using the
// passed logFunc function. Each line is prefixed with tag.
//
// logFunc should have a signature matching log.Println or fmt.Println.
// It is optional. If you omit it, uses fmt.Println for output.
//
// like log.Println: func(...interface{})
//
// like fmt.Println: func(...interface{}) (int, error)
//
func (ob *dataItem) LogStats(tag string, logFunc ...interface{}) {
	log := func(v ...interface{}) { _, _ = fmt.Println(v...) }
	if len(logFunc) > 0 {
		switch fn := logFunc[0].(type) {
		case func(...interface{}): // like log.Println
			log = fn
		case func(...interface{}) (int, error): // like fmt.Println
			log = func(v ...interface{}) { _, _ = fn(v...) }
		}
	}
	log(tag+" name:", ob.Name)
	log(tag+" hash:", fmt.Sprintf("%X", ob.Hash))
	log(tag+" pcs.:", len(ob.CompressedPieces))
	log(tag+" comp:", ob.CompressedSizeInfo, "bytes")
	log(tag+" size:", ob.UncompressedSizeInfo, "bytes")
} //                                                                    LogStats

// Reset discards the contents of the data item and clears its name and hash.
func (ob *dataItem) Reset() {
	ob.Name = ""
	ob.Hash = nil
	ob.CompressedPieces = nil
	ob.CompressedSizeInfo = 0
	ob.UncompressedSizeInfo = 0
} //                                                                       Reset

// Retain changes the Name, Hash, and empties CompressedPieces when the passed
// name, hash and packetCount don't match their current values in the object.
func (ob *dataItem) Retain(name string, hash []byte, packetCount int) {
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

// UnpackBytes joins CompressedPieces and uncompresses
// the resulting bytes to get the original data item.
func (ob *dataItem) UnpackBytes(compressor Compression) ([]byte, error) {
	//
	// join pieces (provided all have been collected) to get compressed data
	if !ob.IsLoaded() {
		return nil, makeError(0xE76AF5, "data item is incomplete")
	}
	compressed := bytes.Join(ob.CompressedPieces, nil)
	ob.CompressedSizeInfo = len(compressed)
	//
	// uncompress data
	ret, err := compressor.Uncompress(compressed)
	if err != nil {
		return nil, makeError(0xE95DFB, err)
	}
	ob.UncompressedSizeInfo = len(ret)
	//
	// hash of uncompressed data should match original hash
	hash, err := getHash(ret)
	if err != nil {
		return nil, makeError(0xE8D61E, err)
	}
	if !bytes.Equal(hash, ob.Hash) {
		return nil, makeError(0xE87D89, "checksum mismatch")
	}
	return ret, nil
} //                                                                 UnpackBytes

// end
