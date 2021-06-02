// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                 /[data_item_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// # Property

// (di *dataItem) IsLoaded() bool
//
// go test -run Test_dataItem_IsLoaded_
//
func Test_dataItem_IsLoaded_(t *testing.T) {
	var d0 dataItem
	if d0.IsLoaded() != false {
		t.Error("0xE3FA57")
	}
	var d1 dataItem
	d1.CompressedPieces = [][]byte{{}, {1}}
	if d1.IsLoaded() != false {
		t.Error("0xE60B69")
	}
	var d2 dataItem
	d2.CompressedPieces = [][]byte{{1}, {23}}
	if d2.IsLoaded() != true {
		t.Error("0xE25AC1")
	}
}

// -----------------------------------------------------------------------------
// # Methods

// (di *dataItem) LogStats(tag string, w io.Writer)
//
// go test -run Test_dataItem_LogStats_
//
func Test_dataItem_LogStats_(t *testing.T) {
	var di = dataItem{
		Key:                  "ItemName",
		Hash:                 []byte{1, 2, 3, 4, 5},
		CompressedPieces:     [][]byte{{6}, {7, 8}, {9, 10, 11}},
		CompressedSizeInfo:   20,
		UncompressedSizeInfo: 50,
	}
	var tlog strings.Builder
	di.LogStats("xyz", &tlog)
	//
	want := "" +
		"xyz  key: ItemName\n" +
		"xyz hash: 0102030405\n" +
		"xyz pcs.: 3\n" +
		"xyz comp: 20 bytes\n" +
		"xyz size: 50 bytes\n"
	//
	got := tlog.String()
	if got != want {
		t.Error("0xE85AA7",
			"\n"+"want:\n", want,
			"\n"+" got:\n", got)
		// fmt.Println("want bytes:", []byte(want))
		// fmt.Println(" got bytes:", []byte(got))
	}
}

// (di *dataItem) Reset()
//
// go test -run Test_dataItem_Reset_
//
func Test_dataItem_Reset_(t *testing.T) {
	var di = dataItem{
		Key:                  "ItemName",
		Hash:                 []byte{1, 2, 3, 4, 5},
		CompressedPieces:     [][]byte{{6}, {7, 8}, {9, 10, 11}},
		CompressedSizeInfo:   20,
		UncompressedSizeInfo: 50,
	}
	di.Reset()
	if di.Key != "" {
		t.Error("0xEA8B3D", "Key not reset")
	}
	if di.Hash != nil {
		t.Error("0xEEA4C6", "Hash not reset")
	}
	if di.CompressedPieces != nil {
		t.Error("0xE3BCE2", "CompressedPieces not reset")
	}
	if di.CompressedSizeInfo != 0 {
		t.Error("0xE04C47", "CompressedSizeInfo not reset")
	}
	if di.UncompressedSizeInfo != 0 {
		t.Error("0xE22CD6", "UncompressedSizeInfo not reset")
	}
}

// (di *dataItem) Retain(k string, hash []byte, packetCount int)
//
// go test -run Test_dataItem_Retain_
//
func Test_dataItem_Retain_(t *testing.T) {
	initDataItem := func() dataItem {
		return dataItem{
			Key:                  "ItemName",
			Hash:                 []byte{1, 2, 3},
			CompressedPieces:     [][]byte{{6}, {7, 8}},
			CompressedSizeInfo:   20,
			UncompressedSizeInfo: 50,
		}
	}
	test := func(k string, hash []byte, packetCount int, want dataItem) {
		var di = initDataItem()
		di.Retain(k, hash, packetCount)
		str := func(di dataItem) string {
			ret := fmt.Sprintf("%#v", di)
			ret = strings.ReplaceAll(ret, "[]uint8", "")
			ret = strings.ReplaceAll(ret, "udpt.dataItem", "")
			return ret
		}
		if !reflect.DeepEqual(di, want) {
			t.Error("0xED7D54", "\n",
				"want:", str(want), "\n",
				" got:", str(di))
		}
	}
	// nothing changed
	test("ItemName", []byte{1, 2, 3}, 2, initDataItem())
	//
	// 'k' parameter changed
	want := dataItem{
		Key:                  "DiffName",
		Hash:                 []byte{1, 2, 3},
		CompressedPieces:     [][]byte{nil, nil},
		CompressedSizeInfo:   0,
		UncompressedSizeInfo: 0,
	}
	test("DiffName", []byte{1, 2, 3}, 2, want)
	//
	// 'hash' parameter changed
	want = dataItem{
		Key:                  "ItemName",
		Hash:                 []byte{6, 7, 8},
		CompressedPieces:     [][]byte{nil, nil},
		CompressedSizeInfo:   0,
		UncompressedSizeInfo: 0,
	}
	test("ItemName", []byte{6, 7, 8}, 2, want)
	//
	// 'packetCount' parameter changed
	want = dataItem{
		Key:                  "ItemName",
		Hash:                 []byte{1, 2, 3},
		CompressedPieces:     [][]byte{nil},
		CompressedSizeInfo:   0,
		UncompressedSizeInfo: 0,
	}
	test("ItemName", []byte{1, 2, 3}, 1, want)
	//
	// all 3 parameters changed
	want = dataItem{
		Key:                  "OtherName",
		Hash:                 []byte{4, 5, 6},
		CompressedPieces:     [][]byte{nil, nil, nil},
		CompressedSizeInfo:   0,
		UncompressedSizeInfo: 0,
	}
	test("OtherName", []byte{4, 5, 6}, 3, want)
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (di *dataItem) UnpackBytes(compressor Compression) ([]byte, error)
//
// go test -run Test_dataItem_UnpackBytes_*

// must succeed
func Test_dataItem_UnpackBytes_1(t *testing.T) {
	source := []byte(strings.Repeat(
		"The quick brown fox jumps over the lazy dog. ", 300,
	))
	zc := &zlibCompressor{}
	comp, err := zc.Compress(source)
	if err != nil {
		t.Error("0xE4A56C", "Compress failed")
	}
	hash := getHash(source)
	var compPieces [][]byte
	{
		a := comp[:]
		for len(a) > 0 {
			n := len(a)
			if n > 50 {
				n = 50
			}
			compPieces = append(compPieces, a[:n])
			a = a[n:]
		}
	}
	var di = dataItem{Hash: hash, CompressedPieces: compPieces}
	// ------------------------------
	uncomp, err := di.UnpackBytes(zc)
	// ------------------------------
	if err != nil {
		t.Error("0xEF6D12", err)
	}
	if !bytes.Equal(source, uncomp) {
		t.Error("0xE91A65", "corrupted data")
	}
	if !bytes.Equal(hash, di.Hash) {
		t.Error("0xEC4E68", "corrupted hash")
	}
	if di.CompressedSizeInfo != len(comp) {
		t.Error("0xEB4A34", "wrong CompressedSizeInfo")
	}
	if di.UncompressedSizeInfo != len(source) {
		t.Error("0xEC1E61", "wrong UncompressedSizeInfo")
	}
}

// must fail trying to unpack an empty item
func Test_dataItem_UnpackBytes_2(t *testing.T) {
	zc := &zlibCompressor{}
	var di0 dataItem
	data, err := di0.UnpackBytes(zc)
	if data != nil {
		t.Error("0xED52E6")
	}
	if !matchError(err, "data item is incomplete") {
		t.Error("0xEE0C63", "wrong error:", err)
	}
}

// must fail when everything succeeds but the hash is wrong
func Test_dataItem_UnpackBytes_3(t *testing.T) {
	source := []byte(strings.Repeat(
		"The quick brown fox jumps over the lazy dog. ", 300,
	))
	zc := &zlibCompressor{}
	comp, err := zc.Compress(source)
	if err != nil {
		t.Error("0xE70C74", "Compress failed")
	}
	var compPieces [][]byte
	{
		a := comp[:]
		for len(a) > 0 {
			n := len(a)
			if n > 50 {
				n = 50
			}
			compPieces = append(compPieces, a[:n])
			a = a[n:]
		}
	}
	var di = dataItem{Hash: getHash(source), CompressedPieces: compPieces}
	zc = &zlibCompressor{}
	// ------------------------------
	di.Hash = []byte{0} // <- this must cause it to fail
	uncomp, err := di.UnpackBytes(zc)
	// ------------------------------
	if uncomp != nil {
		t.Error("0xED14FA")
	}
	if !matchError(err, "hash mismatch") {
		t.Error("0xEA19E1", "wrong error:", err)
	}
}

// must fail to uncompress an item containing garbage bytes
func Test_dataItem_UnpackBytes_4(t *testing.T) {
	var di = dataItem{
		Hash: []byte{0xA1, 0x96, 0x9E, 0xBF, 0x93, 0xE5},
		CompressedPieces: [][]byte{{
			0xC6, 0x44, 0x0D, 0xAC, 0xA9, 0x55, 0x4D, 0xEF,
			0xA1, 0x93, 0x8D, 0x41, 0x80, 0x61, 0x29, 0xC2,
		}},
	}
	zc := &zlibCompressor{}
	uncomp, err := di.UnpackBytes(zc)
	if uncomp != nil {
		t.Error("0xE59B01")
	}
	if !matchError(err, "zlib") {
		t.Error("0xEF8DE2", "wrong error:", err)
	}
}

// end
