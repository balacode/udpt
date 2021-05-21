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

// (ob *dataItem) IsLoaded() bool
//
// go test -run Test_dataItem_IsLoaded_
//
func Test_dataItem_IsLoaded_(t *testing.T) {
	var dataItem0 dataItem
	if dataItem0.IsLoaded() != false {
		t.Error("0xE3FA57 dataItem0.IsLoaded() expect:false got:true")
	}
	var dataItem1 dataItem
	dataItem1.CompressedPieces = [][]byte{{}, {1}}
	if dataItem1.IsLoaded() != false {
		t.Error("0xE60B69 dataItem1.IsLoaded() expect:false got:true")
	}
	var dataItem2 dataItem
	dataItem2.CompressedPieces = [][]byte{{1}, {23}}
	if dataItem2.IsLoaded() != true {
		t.Error("0xE25AC1 dataItem2.IsLoaded() expect:true got:false")
	}
} //                                                     Test_dataItem_IsLoaded_

// -----------------------------------------------------------------------------
// # Methods

// (ob *dataItem) LogStats(tag string, w io.Writer)
//
// go test -run Test_dataItem_LogStats_
//
func Test_dataItem_LogStats_(t *testing.T) {
	var sb strings.Builder
	fmtPrintln := func(v ...interface{}) (int, error) {
		sb.WriteString(fmt.Sprintln(v...))
		return 0, nil
	}
	logPrintln := func(v ...interface{}) {
		sb.WriteString(fmt.Sprintln(v...))
	}
	test := func(logFunc interface{}) {
		var di = dataItem{
			Name:                 "ItemName",
			Hash:                 []byte{1, 2, 3, 4, 5},
			CompressedPieces:     [][]byte{{6}, {7, 8}, {9, 10, 11}},
			CompressedSizeInfo:   20,
			UncompressedSizeInfo: 50,
		}
		sb.Reset()
		di.LogStats("xyz", &sb)
		got := sb.String()
		//
		expect := "" +
			"xyz name: ItemName\n" +
			"xyz hash: 0102030405\n" +
			"xyz pcs.: 3\n" +
			"xyz comp: 20 bytes\n" +
			"xyz size: 50 bytes\n"
		//
		if got != expect {
			t.Error("\n"+"expect:\n", expect, "\n"+"got:\n", got)
			fmt.Println([]byte(expect))
			fmt.Println([]byte(got))
		}
	}
	test(fmtPrintln)
	test(logPrintln)
} //                                                     Test_dataItem_LogStats_

// (ob *dataItem) Reset()
//
// go test -run Test_dataItem_Reset_
//
func Test_dataItem_Reset_(t *testing.T) {
	var di = dataItem{
		Name:                 "ItemName",
		Hash:                 []byte{1, 2, 3, 4, 5},
		CompressedPieces:     [][]byte{{6}, {7, 8}, {9, 10, 11}},
		CompressedSizeInfo:   20,
		UncompressedSizeInfo: 50,
	}
	di.Reset()
	if di.Name != "" {
		t.Error("0xEA8B3D Name not reset")
	}
	if di.Hash != nil {
		t.Error("0xEEA4C6 Hash not reset")
	}
	if di.CompressedPieces != nil {
		t.Error("0xE3BCE2 CompressedPieces not reset")
	}
	if di.CompressedSizeInfo != 0 {
		t.Error("0xE04C47 CompressedSizeInfo not reset")
	}
	if di.UncompressedSizeInfo != 0 {
		t.Error("0xE22CD6 UncompressedSizeInfo not reset")
	}
} //                                                        Test_dataItem_Reset_

// (ob *dataItem) Retain(name string, hash []byte, packetCount int)
//
// go test -run Test_dataItem_Retain_
//
func Test_dataItem_Retain_(t *testing.T) {
	initDataItem := func() dataItem {
		return dataItem{
			Name:                 "ItemName",
			Hash:                 []byte{1, 2, 3},
			CompressedPieces:     [][]byte{{6}, {7, 8}},
			CompressedSizeInfo:   20,
			UncompressedSizeInfo: 50,
		}
	}
	test := func(name string, hash []byte, packetCount int, expect dataItem) {
		var di = initDataItem()
		di.Retain(name, hash, packetCount)
		str := func(di dataItem) string {
			ret := fmt.Sprintf("%#v", di)
			ret = strings.ReplaceAll(ret, "[]uint8", "")
			ret = strings.ReplaceAll(ret, "udpt.dataItem", "")
			return ret
		}
		if !reflect.DeepEqual(di, expect) {
			t.Error("0xED7D54 Retain() failed\n",
				"expect:", str(expect), "\n",
				"   got:", str(di))
		}
	}
	{
		// nothing changed
		test("ItemName", []byte{1, 2, 3}, 2, initDataItem())
	}
	{
		// 'name' parameter changed
		expect := dataItem{
			Name:                 "DiffName",
			Hash:                 []byte{1, 2, 3},
			CompressedPieces:     [][]byte{nil, nil},
			CompressedSizeInfo:   0,
			UncompressedSizeInfo: 0,
		}
		test("DiffName", []byte{1, 2, 3}, 2, expect)
	}
	{
		// 'hash' parameter changed
		expect := dataItem{
			Name:                 "ItemName",
			Hash:                 []byte{6, 7, 8},
			CompressedPieces:     [][]byte{nil, nil},
			CompressedSizeInfo:   0,
			UncompressedSizeInfo: 0,
		}
		test("ItemName", []byte{6, 7, 8}, 2, expect)
	}
	{
		// 'packetCount' parameter changed
		expect := dataItem{
			Name:                 "ItemName",
			Hash:                 []byte{1, 2, 3},
			CompressedPieces:     [][]byte{nil},
			CompressedSizeInfo:   0,
			UncompressedSizeInfo: 0,
		}
		test("ItemName", []byte{1, 2, 3}, 1, expect)
	}
	{
		// all 3 parameters changed
		expect := dataItem{
			Name:                 "OtherName",
			Hash:                 []byte{4, 5, 6},
			CompressedPieces:     [][]byte{nil, nil, nil},
			CompressedSizeInfo:   0,
			UncompressedSizeInfo: 0,
		}
		test("OtherName", []byte{4, 5, 6}, 3, expect)
	}
} //                                                       Test_dataItem_Retain_

// (ob *dataItem) UnpackBytes(compressor Compression) ([]byte, error)
//
// go test -run Test_dataItem_UnpackBytes_
//
func Test_dataItem_UnpackBytes_(t *testing.T) {
	zc := &zlibCompressor{}
	{
		var dataItem0 dataItem
		data, err := dataItem0.UnpackBytes(zc)
		if data != nil {
			t.Error("0xED52E6 dataItem0.UnpackBytes()",
				"returned: data != nil, expect: data == nil")
		}
		if err == nil {
			t.Error("0xEE0C63 dataItem0.UnpackBytes()",
				"returned: error == nil, expect: error != nil")
		}
	}
	{
		source := []byte(strings.Repeat(
			"The quick brown fox jumps over the lazy dog. ", 300,
		))
		hash := getHash(source)
		compressed, err := zc.Compress(source)
		if err != nil {
			t.Error("0xE70C74 Compress failed")
		}
		var compPieces [][]byte
		{
			a := compressed[:]
			for len(a) > 0 {
				n := len(a)
				if n > 50 {
					n = 50
				}
				compPieces = append(compPieces, a[:n])
				a = a[n:]
			}
		}
		var dataItem1 = dataItem{
			Hash:             hash,
			CompressedPieces: compPieces,
		}
		uncompressed, err := dataItem1.UnpackBytes(zc)
		if err != nil {
			t.Error("0xEF6D12 UnpackBytes:", err)
		}
		if !bytes.Equal(source, uncompressed) {
			t.Error("0xE91A65 UnpackBytes: corrupted data")
		}
		if !bytes.Equal(hash, dataItem1.Hash) {
			t.Error("0xEC4E68 UnpackBytes: corrupted hash")
		}
		if dataItem1.CompressedSizeInfo != len(compressed) {
			t.Error("0xEB4A34",
				"CompressedSizeInfo", dataItem1.CompressedSizeInfo,
				"!= len(compressed)", len(compressed))
		}
		if dataItem1.UncompressedSizeInfo != len(source) {
			t.Error("0xEC1E61",
				"UncompressedSizeInfo", dataItem1.UncompressedSizeInfo,
				"!= len(source)", len(source))
		}
		{
			dataItem1.Hash = []byte{0}
			uncompressed, err := dataItem1.UnpackBytes(zc)
			if uncompressed != nil {
				t.Error("0xED14FA")
			}
			if !matchError(err, "hash mismatch") {
				t.Error("0xEA19E1")
			}
		}
	}
	{
		// try to uncompress an item containing garbage bytes
		var dataItem3 = dataItem{
			Hash: []byte{0xA1, 0x96, 0x9E, 0xBF, 0x93, 0xE5},
			CompressedPieces: [][]byte{{
				0xC6, 0x44, 0x0D, 0xAC, 0xA9, 0x55, 0x4D, 0xEF,
				0xA1, 0x93, 0x8D, 0x41, 0x80, 0x61, 0x29, 0xC2,
			}},
		}
		uncompressed, err := dataItem3.UnpackBytes(zc)
		if !matchError(err, "zlib") {
			t.Error("0xEF8DE2")
		}
		if uncompressed != nil {
			t.Error("0xE59B01")
		}
	}
} //                                                  Test_dataItem_UnpackBytes_

// end
